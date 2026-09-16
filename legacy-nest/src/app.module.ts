import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { ContactsModule } from './contacts/contacts.module';

@Module({
  imports: [
    // SQLite keeps this repo dependency-free to run. The real GloveBoxCRM
    // legacy service runs Postgres — swap this for the same connection
    // string go-service uses if you want to run both against one DB.
    TypeOrmModule.forRoot({
      type: 'sqlite',
      database: 'legacy.sqlite',
      autoLoadEntities: true,
      synchronize: true, // fine for this demo; never in real prod
    }),
    ContactsModule,
  ],
})
export class AppModule {}
