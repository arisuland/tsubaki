-- AlterTable
ALTER TABLE "users" ADD COLUMN     "avatarUrl" TEXT,
ADD COLUMN     "gravatarEmail" TEXT,
ADD COLUMN     "useGravatar" BOOLEAN NOT NULL DEFAULT false;
