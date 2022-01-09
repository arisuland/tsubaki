-- CreateEnum
CREATE TYPE "AccessTokenScope" AS ENUM ('PUBLIC_WRITE', 'REPO_CREATE', 'REPO_DELETE', 'REPO_UPDATE');

-- AlterTable
ALTER TABLE "access_tokens" ADD COLUMN     "expiresIn" TIMESTAMP(3),
ADD COLUMN     "scopes" "AccessTokenScope"[];

-- CreateTable
CREATE TABLE "subprojects" (
    "description" TEXT,
    "updated_at" TIMESTAMP(3) NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "parentId" TEXT NOT NULL,
    "name" TEXT NOT NULL,
    "id" TEXT NOT NULL,

    CONSTRAINT "subprojects_pkey" PRIMARY KEY ("id")
);

-- AddForeignKey
ALTER TABLE "subprojects" ADD CONSTRAINT "subprojects_parentId_fkey" FOREIGN KEY ("parentId") REFERENCES "projects"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
