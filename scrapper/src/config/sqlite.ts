import { Database } from 'sqlite3';

export async function createDbConnection(): Promise<Database> {
  return new Promise((resolve, reject) => {
    const db = new Database('./scraped_data.db', (err) => {
      if (err) {
        reject(err);
      } else {
        resolve(db);
      }
    });
  });
}

