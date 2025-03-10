import { scraperService } from "./services/scrapper";
import { ApartmentRepository } from "./repositories/apartment.repository";
import { createDbConnection } from "./config/sqlite";
import logger from "./config/logger";
import sleep from "./utils/sleep";

const scrapeFullList = async () => {
  const db = await createDbConnection();
  const apartmentRepo = new ApartmentRepository(db);
  await apartmentRepo.initialize();

  try {
    let firstId: string | null = null;
    let page = 1;


    while(true) {
      const scrapedList = await scraperService.list(page)
      if (firstId && scrapedList.propertyIds.includes(firstId)) break;

      if (!firstId && scrapedList.propertyIds.length > 0) {
        firstId = scrapedList.propertyIds[0];
        logger.debug(`FirstId: ${firstId}`)
      }

      logger.info(`Analizing page: ${page}`)

      for (const apartmentId of scrapedList.propertyIds){
        if(await apartmentRepo.exists(apartmentId)){
          logger.info(`Already saved: ${apartmentId}`)
          continue;
        }

        await sleep(Math.floor(Math.random() * 5000) + 5000)
        const scrapedData = await scraperService.apartment(apartmentId);

        await apartmentRepo.save({
          id: apartmentId,
          title: scrapedData.title,
          details: JSON.stringify(scrapedData.details),
          pricePerSqm: scrapedData.priceFeatures.pricePerSqm,
          propertyPrice: scrapedData.priceFeatures.propertyPrice,
          communityFees: scrapedData.priceFeatures.communityFees,
          location: scrapedData.location.toString(),
          description: scrapedData.description,
          createdAt: new Date()
        });
      }

      page++
      await sleep(60000)
    }
  } catch (error) {
    apartmentRepo.close()
    console.error(error instanceof Error ? error : error);
  }
}

const runScraper = async (): Promise<void> => {
    await scrapeFullList()
}

runScraper();
