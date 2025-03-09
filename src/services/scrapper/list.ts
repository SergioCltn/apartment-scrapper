import axios, { AxiosResponse } from 'axios';
import * as cheerio from 'cheerio';

interface ScrapedData {
  propertyIds: string[]
}

export default async function scrapeList(page: number): Promise<ScrapedData> {
  const headers = {
    userAgent : 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:135.0) Gecko/20100101 Firefox/135.0',
    acceptLanguage : 'en-US,en;q=0.5',
  }

  try {
    const response: AxiosResponse<string> = await axios.get<string>(`https://www.idealista.com/areas/venta-viviendas/con-precio-hasta_190000,precio-desde_80000/${'pagina-'+page}?shape=%28%28%28%7Bj%7C%7EE%7EhfZ%7CH_%40%7EEaAdEeDtL%7BLjDwElFyIxDyHlAwFmAgHgEkJqCuEdCnO%7BFsCwNwFcVwGuLuDyFsBiGsCqJgGiN%7DKcHwGwIiJoHiGsEsCoCqAoYqPwPkKaFuD%7DHgHoOoOwBuCoH%7DMoAcDoH_OgGaQoAuEmFaPYaBcCgGqCuDqCsC_FuCyDaBeEq%40aF_%40cC%3FeEn%40sEbBmFrCiBpA%7DHhHmA%60BcH%7CLiBvEmAtE_DhIoAfGaAhIe%40%7ENXfHr%40xHpCnNxD%7ENzHdUvBrDxDzIlAtDdEhJvBvFzAbCxDhJxK%7CLlM%60Q%60FhHjDbDbCtCfGtFfEtCzAn%40lRdEdErApCN%7EJrBjK%3FrEN%60McCpCqAlFcBpC_%40dE%3FzFpAzTzJvIbDnMxGzMdFxKbChGp%40%60H%5EpH%3F%29%29%29&ordenado-por=fecha-publicacion-desc`, {
      headers: {
        'User-Agent': headers.userAgent,
        'Accept-Language': headers.acceptLanguage,
      },
      timeout: 10000,
    });

    const $ = cheerio.load(response.data);

    const propertyIds: string[] = [];
    $('.items-container.items-list article.item').each((_, element) => {
      const id = $(element).attr('data-element-id');
      if (id && !propertyIds.includes(id)) {
        propertyIds.push(id);
      }
    });

    return { propertyIds };
  } catch (error) {
    const errorMessage = error instanceof Error ? error.message : 'Unknown error';
    throw new Error(`Failed to scrape Idealista: ${errorMessage}`);
  }
}
