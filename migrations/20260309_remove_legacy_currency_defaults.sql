ALTER TABLE `purchase_order`
  ALTER COLUMN `currency` DROP DEFAULT;

ALTER TABLE `purchase_order_item`
  ALTER COLUMN `currency` DROP DEFAULT;

ALTER TABLE `shipment`
  ALTER COLUMN `currency` DROP DEFAULT;

ALTER TABLE `shipment_item`
  ALTER COLUMN `currency` DROP DEFAULT;

ALTER TABLE `packaging_item`
  ALTER COLUMN `currency` DROP DEFAULT;

ALTER TABLE `shipping_rate`
  ALTER COLUMN `currency` DROP DEFAULT;
