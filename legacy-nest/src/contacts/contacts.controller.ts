import { Body, Controller, Delete, Get, Param, Patch, Post } from '@nestjs/common';
import { ContactsService } from './contacts.service';
import { Policyholder } from './policyholder.entity';

// Route shape intentionally matches go-service's /api/v1/policyholders
// endpoints, so the gateway can strangler-route between the two without
// clients noticing a contract change.
@Controller('api/v1/policyholders')
export class ContactsController {
  constructor(private readonly contacts: ContactsService) {}

  @Post()
  create(@Body() body: Partial<Policyholder>) {
    return this.contacts.create(body);
  }

  @Get()
  findAll() {
    return this.contacts.findAll();
  }

  @Get(':id')
  findOne(@Param('id') id: string) {
    return this.contacts.findOne(id);
  }

  @Patch(':id')
  update(@Param('id') id: string, @Body() body: Partial<Policyholder>) {
    return this.contacts.update(id, body);
  }

  @Delete(':id')
  remove(@Param('id') id: string) {
    return this.contacts.remove(id);
  }
}
