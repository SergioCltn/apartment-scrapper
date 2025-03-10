import axios, { AxiosResponse } from 'axios';
import * as cheerio from 'cheerio';

interface PriceFeatures {
  propertyPrice?: string,
  pricePerSqm?: string, 
  communityFees?: string,
}

type Details = Record<string, string[]>;

interface ScrapedData {
  title: string; 
  description: string; 
  details: Details;
  detailInfoTag: string;
  titleMinor: string;     
  infoFeatures: string[];
  priceFeatures: PriceFeatures;
  location: string[];   
}

const userAgents = [
  'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36',
  'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.0 Safari/605.1.15',
  'Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/115.0',
];

function getRandomUserAgent(): string {
  return userAgents[Math.floor(Math.random() * userAgents.length)];
}

export default async function scrapeApartment(apartmentId: string): Promise<ScrapedData> {
  const headers = {
    userAgent : 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:135.0) Gecko/20100101 Firefox/135.0',
    acceptLanguage : 'en-US,en;q=0.5',
  }

  try {
    const response: AxiosResponse<string> = await axios.get<string>(`https://www.idealista.com/inmueble/${apartmentId}`, {
      headers: {
        'User-Agent': getRandomUserAgent(),
        'Accept-Language': headers.acceptLanguage,
      },
      timeout: 10000,
    });

    const $ = cheerio.load(response.data);

    const description: string = $('.comment .adCommentsLanguage p').first().text().trim() || 'No comment found';

    const title: string = $('.main-info__title-main').first().text().trim() || 'No main title found';

    const titleMinor: string = $('.main-info__title-minor').first().text().trim() || 'No minor title found';

    const detailInfoTag: string = $('.detail-info-tags .tag').first().text().trim() || 'No minor title found';

    const infoFeatures: string[] = [];
    $('.info-features span').each((_, element) => {
      const text = $(element).text().trim();
      if (text) infoFeatures.push(text);
    });

    const details: Details = {};
    $('#details .details-property-h2').each((_, element) => {
      const sectionTitle = $(element).text().trim();
      if (!sectionTitle) return;

      const ul = $(element).next('.details-property_features').find('ul');
      const items: string[] = [];
      ul.find('li').each((_, li) => {
        const text = $(li).text().trim();
        if (text) items.push(text);
      });

      if (items.length > 0) details[sectionTitle] = items;
    });

    const location: string[] = [];
    $('#headerMap ul li.header-map-list').each((_, element) => {
      const text = $(element).text().trim();
      if (text) location.push(text);
    });

    const priceFeatures: PriceFeatures = {};

    $('.price-features__container p').each((_, element) => {
      const $element = $(element);
      const spans = $element.find('.flex-feature-details');

      if (spans.length === 2) {
        const label = spans.eq(0).text().trim();
        if (!label) return;

        if (spans.length === 2) {
          const value = spans.eq(1).text().trim();

          if (label.includes('Precio del inmueble')) {
            priceFeatures.propertyPrice = value || undefined;
          } else if (label.includes('Precio por m²')) {
            priceFeatures.pricePerSqm = value || undefined;
          }
        }
      } else if (spans.length === 1) {
        const text = spans.eq(0).text().trim();
        if (text.startsWith('Gastos de comunidad')) {
          priceFeatures.communityFees = text.replace('Gastos de comunidad', '').trim() || undefined;
        }
      }
    });

    return { title, titleMinor, detailInfoTag, description, infoFeatures, details, priceFeatures, location  };
  } catch (error) {
    const errorMessage = error instanceof Error ? error.message : 'Unknown error';
    throw new Error(`Failed to scrape Idealista: ${errorMessage}`);
  }
}
