import { Injectable, NotFoundException } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { Policyholder } from './policyholder.entity';

@Injectable()
export class ContactsService {
  constructor(
    @InjectRepository(Policyholder)
    private readonly repo: Repository<Policyholder>,
  ) {}

  create(data: Partial<Policyholder>) {
    const entity = this.repo.create(data);
    return this.repo.save(entity);
  }

  findAll() {
    return this.repo.find({ order: { createdAt: 'DESC' } });
  }

  async findOne(id: string) {
    const found = await this.repo.findOneBy({ id });
    if (!found) throw new NotFoundException('policyholder not found');
    return found;
  }

  async update(id: string, data: Partial<Policyholder>) {
    await this.findOne(id);
    await this.repo.update(id, data);
    return this.findOne(id);
  }

  async remove(id: string) {
    await this.findOne(id);
    await this.repo.delete(id);
  }
}
