import { Database } from 'sqlite3';
import { Apartment } from '../models/aparments';
import logger from '../config/logger';

export class ApartmentRepository {
  private db: Database;

  constructor(db: Database) {
    this.db = db;
  }

  async initialize(): Promise<void> {
    return new Promise((resolve, reject) => {
      this.db.run(`
        CREATE TABLE IF NOT EXISTS apartments (
          id INTEGER PRIMARY KEY ,
          title TEXT,
          details TEXT,
          propertyPrice TEXT,
          pricePerSqm TEXT,
          communityFees TEXT,
          location TEXT,
          description TEXT,
          createdAt TEXT NOT NULL,
          updatedAt TEXT
        )
      `, (err) => {
        if (err) reject(err);
        else {
              logger.info(`Initializing database connection`)
              resolve();
            }
      });
    });
  }

  async save(apartment: Apartment): Promise<void> {
    const query = `
      INSERT INTO apartments (
        id, title, location, description, createdAt, updatedAt, propertyPrice, pricePerSqm, communityFees, details
      ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `;
    
    return new Promise((resolve, reject) => {
      this.db.run(
        query,
        [
          apartment.id,
          apartment.title,
          apartment.location, 
          apartment.description,
          apartment.createdAt.toISOString(),
          apartment.updatedAt?.toISOString(),
          apartment.propertyPrice,
          apartment.pricePerSqm,
          apartment.communityFees,
          apartment.details,
        ],
        (err) => {
          if (err) {
            reject(new Error(`Failed to save apartment ${apartment.id}: ${err.message}`));
          } else {
            logger.info(`Saved apartment with id: ${apartment.id}`)
            resolve();
          }
        }
      );
    });
  }

  async exists(id: string): Promise<boolean> {
    return new Promise((resolve, reject) => {
      this.db.get(
        'SELECT 1 FROM apartments WHERE id = ? LIMIT 1',
        [id],
        (err, row) => {
          if (err) reject(err);
          else resolve(!!row);
        }
      );
    });
  }

  async findById(id: string): Promise<Apartment | null> {
    return new Promise((resolve, reject) => {
      this.db.get(
        'SELECT * FROM apartments WHERE id = ?',
        [id],
        (err, row: Apartment) => {
          if (err) {
            reject(err);
          } else if (!row) {
            resolve(null);
          } else {
            resolve({
              id: row.id,
              title: row.title,
              communityFees: row.communityFees,
              propertyPrice: row.propertyPrice,
              pricePerSqm: row.pricePerSqm,
              details: row.details,
              location: row.location,
              description: row.description,
              createdAt: new Date(row.createdAt),
              updatedAt: row.updatedAt ? new Date(row.updatedAt) : undefined,
            });
          }
        }
      );
    });
  }

  async update(id: string, data: Partial<Apartment>): Promise<void> {
    const fields: string[] = [];
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const values: any[] = [];

    if (data.title !== undefined) {
      fields.push('title = ?');
      values.push(data.title);
    }
    if (data.pricePerSqm !== undefined) {
      fields.push('pricePerSqm = ?');
      values.push(data.pricePerSqm);
    }
    if (data.propertyPrice !== undefined) {
      fields.push('propertyPrice = ?');
      values.push(data.propertyPrice);
    }
    if (data.communityFees !== undefined) {
      fields.push('communityFees = ?');
      values.push(data.communityFees);
    }
    if (data.location !== undefined) {
      fields.push('location = ?');
      values.push(data.location);
    }
    if (data.description !== undefined) {
      fields.push('description = ?');
      values.push(data.description);
    }
    if (data.details !== undefined) {
      fields.push('details = ?');
      values.push(data.details);
    }

    fields.push('updatedAt = ?');
    values.push(new Date().toISOString());
    values.push(id);

    if (fields.length === 0) return;

    const query = `UPDATE apartments SET ${fields.join(', ')} WHERE id = ?`;
    
    return new Promise((resolve, reject) => {
      this.db.run(query, values, (err) => {
        if (err) reject(err);
        else {
            logger.info(`Updated apartment with id: ${id}`)
            resolve()
          };
      });
    });
  }

  async close(): Promise<void> {
    return new Promise((resolve, reject) => {
      this.db.close((err) => {
        if (err) reject(err);
        else { 
            logger.info(`Closing database connection`)
            resolve(); 
          }
      });
    });
  }
}
