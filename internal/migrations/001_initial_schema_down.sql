-- +migrate Down
DO $$
BEGIN
    DROP TRIGGER IF EXISTS update_vehicles_updated_at ON vehicles;
    DROP TRIGGER IF EXISTS update_routes_updated_at ON routes;
    DROP TRIGGER IF EXISTS update_deliveries_updated_at ON deliveries;
    DROP TRIGGER IF EXISTS update_inventory_updated_at ON inventory;
    DROP TRIGGER IF EXISTS update_warehouses_updated_at ON warehouses;
    DROP TRIGGER IF EXISTS update_products_updated_at ON products;
END $$;

DROP INDEX IF EXISTS idx_tracking_events_tracking_id;
DROP INDEX IF EXISTS idx_tracking_info_delivery_id;
DROP INDEX IF EXISTS idx_deliveries_status;
DROP INDEX IF EXISTS idx_inventory_movements_product_id;
DROP INDEX IF EXISTS idx_inventory_product_id;
DROP INDEX IF EXISTS idx_products_sku;

DROP TABLE IF EXISTS delivery_trackings;
DROP TABLE IF EXISTS tracking_events;
DROP TABLE IF EXISTS tracking_info;
DROP TABLE IF EXISTS delivery_items;
DROP TABLE IF EXISTS routes;
DROP TABLE IF EXISTS deliveries;
DROP TABLE IF EXISTS inventory_movements;
DROP TABLE IF EXISTS inventory;
DROP TABLE IF EXISTS vehicles;
DROP TABLE IF EXISTS warehouses;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS notifications; 