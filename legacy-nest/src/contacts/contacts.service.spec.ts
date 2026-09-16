import { Test, TestingModule } from '@nestjs/testing';
import { getRepositoryToken } from '@nestjs/typeorm';
import { ContactsService } from './contacts.service';
import { Policyholder } from './policyholder.entity';

describe('ContactsService', () => {
  let service: ContactsService;

  const mockRepo = {
    create: jest.fn((dto) => dto),
    save: jest.fn((dto) => Promise.resolve({ id: 'uuid-1', ...dto })),
    find: jest.fn(() => Promise.resolve([])),
    findOneBy: jest.fn(({ id }) =>
      id === 'uuid-1'
        ? Promise.resolve({ id: 'uuid-1', fullName: 'Jane Doe', email: 'jane@example.com' })
        : Promise.resolve(null),
    ),
    update: jest.fn(() => Promise.resolve()),
    delete: jest.fn(() => Promise.resolve()),
  };

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      providers: [
        ContactsService,
        {
          provide: getRepositoryToken(Policyholder),
          useValue: mockRepo,
        },
      ],
    }).compile();

    service = module.get<ContactsService>(ContactsService);
  });

  it('should be defined', () => {
    expect(service).toBeDefined();
  });

  it('should find all contacts', async () => {
    const result = await service.findAll();
    expect(result).toEqual([]);
    expect(mockRepo.find).toHaveBeenCalled();
  });

  it('should find one contact by id', async () => {
    const result = await service.findOne('uuid-1');
    expect(result.fullName).toBe('Jane Doe');
  });
});
