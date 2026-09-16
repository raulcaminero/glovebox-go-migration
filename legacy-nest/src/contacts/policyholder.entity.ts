import { Column, CreateDateColumn, Entity, PrimaryGeneratedColumn, UpdateDateColumn } from 'typeorm';

// Mirrors the schema in go-service/db/migrations/0001_init.sql. This is the
// entity as it exists today in the legacy service — the migration's job is
// to make internal/domain/policyholder.go the new source of truth for this
// shape while keeping the API contract stable during rollout.
@Entity('policyholders')
export class Policyholder {
  @PrimaryGeneratedColumn('uuid')
  id!: string;

  @Column()
  fullName!: string;

  @Column({ unique: true })
  email!: string;

  @Column({ nullable: true })
  phone?: string;

  @CreateDateColumn()
  createdAt!: Date;

  @UpdateDateColumn()
  updatedAt!: Date;
}
