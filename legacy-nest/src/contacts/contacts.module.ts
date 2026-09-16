import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { Policyholder } from './policyholder.entity';
import { ContactsService } from './contacts.service';
import { ContactsController } from './contacts.controller';

@Module({
  imports: [TypeOrmModule.forFeature([Policyholder])],
  controllers: [ContactsController],
  providers: [ContactsService],
})
export class ContactsModule {}
