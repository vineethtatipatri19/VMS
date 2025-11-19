# Database Migrations

This directory contains SQL migration files that create and modify the database schema.

## Migration Files

The migrations are applied in numerical order:

1. **001_init.sql** - Creates core tables (customers, inventory_items, transactions)
2. **002_users.sql** - Creates users table for authentication
3. **003_add_indexes.sql** - Adds performance indexes
4. **004_enhance_entities.sql** - Expands entity fields (35+ fields per item)
5. **005_expand_units.sql** - Adds support for more unit types
6. **006_add_transaction_status.sql** - Adds transaction status tracking
7. **007_add_soft_delete.sql** - Implements audit trail with soft deletes

## Important Notes

### Automatic Migration Issue

**Current Limitation:** The backend's automatic migration system using golang-migrate has a path resolution issue. The library expects files named `001_init.up.sql` and `001_init.down.sql`, but our files use just `.sql` suffix.

**Why not rename?** We keep the simple `.sql` naming because:
- These migrations are applied manually during setup
- Keeps filenames simple and readable
- Avoids needing separate "down" migrations for each change

### How Migrations Are Applied

**Via setup.sh (Recommended):**
```bash
bash setup.sh
```
The setup script automatically applies all migrations in order before loading demo data.

**Manually:**
```bash
# Apply all migrations
for file in infra/migrations/*.sql; do 
  echo "Applying $(basename $file)..."
  docker exec -i pgvms-postgres psql -U pgvms_user -d pgvms < "$file"
done
```

**Individual migration:**
```bash
docker exec -i pgvms-postgres psql -U pgvms_user -d pgvms < infra/migrations/007_add_soft_delete.sql
```

### Verification

After running migrations, verify all tables exist:

```bash
docker exec pgvms-postgres psql -U pgvms_user -d pgvms -c "\dt"
```

Expected tables:
- customers
- inventory_items
- transactions
- sale_items
- crate_ledger
- wastage_log
- expiry_alerts
- users
- payment_schedules
- pricing_tiers
- price_history

## Migration Order Matters

⚠️ **Important:** Always apply migrations in numerical order. Later migrations depend on tables and columns created by earlier ones.

## Critical Migration: 007_add_soft_delete.sql

This migration is essential for the audit trail system. It adds:
- `deleted_at` TIMESTAMPTZ - Soft delete timestamp
- `deleted_by` TEXT - Who performed the deletion
- `deletion_reason` TEXT - Why the record was deleted

These fields enable the "nothing is truly deleted" feature with complete audit trail.

## Docker Compose Configuration

The `docker-compose.yml` mounts migrations to the backend container:

```yaml
backend:
  volumes:
    - ./infra/migrations:/app/infra/migrations:ro
```

However, migrations should NOT be mounted to the database init directory:

```yaml
# ❌ DON'T DO THIS (causes demo data to run before tables exist):
db:
  volumes:
    - ./infra/local:/docker-entrypoint-initdb.d:ro

# ✅ CORRECT (only persist data):
db:
  volumes:
    - postgres_data:/var/lib/postgresql/data
```

## Adding New Migrations

When adding new migrations:

1. **Name with next number:** `008_your_migration.sql`
2. **Include comments:** Explain what the migration does
3. **Make idempotent if possible:** Use `IF NOT EXISTS` where applicable
4. **Test thoroughly:** Apply to fresh database and verify
5. **Update this README:** Document the new migration

Example:
```sql
-- 008_add_supplier_ratings.sql
-- Adds supplier rating system to inventory items

ALTER TABLE inventory_items 
ADD COLUMN IF NOT EXISTS supplier_rating NUMERIC CHECK (supplier_rating >= 0 AND supplier_rating <= 5);

CREATE INDEX IF NOT EXISTS idx_inventory_supplier_rating 
ON inventory_items(supplier_rating) 
WHERE deleted_at IS NULL;
```

## Production Considerations

For production deployments:

1. **Backup first:** Always backup the database before applying migrations
2. **Test migrations:** Apply to staging environment first
3. **Monitor closely:** Watch for errors during migration
4. **Plan rollback:** Have a rollback plan ready
5. **Low traffic window:** Apply during maintenance window if possible

## Troubleshooting

### "relation does not exist"
Migrations haven't been applied yet. Run setup script or apply manually.

### "column already exists"
Migration was partially applied. Either:
- Drop and recreate database: `docker-compose down -v && bash setup.sh`
- Or manually fix the specific table

### "Migration warning: first .: file does not exist"
This is expected from the backend's golang-migrate. Migrations must be applied manually via the setup script.

## See Also

- [SETUP.md](../../SETUP.md) - Complete setup guide
- [QUICK_REFERENCE.md](../../QUICK_REFERENCE.md) - Common commands
- [infra/local/](../local/) - Demo data SQL files (applied AFTER migrations)
