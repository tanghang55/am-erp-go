UPDATE `shipment`
SET `inventory_locked` = 0
WHERE `status` IN ('SHIPPED', 'DELIVERED', 'CANCELLED')
  AND `inventory_locked` = 1;

UPDATE `shipment` s
LEFT JOIN (
    SELECT `reference_id`, COUNT(*) AS `allocate_count`
    FROM `inventory_movement`
    WHERE `reference_type` = 'SHIPMENT'
      AND `movement_type` = 'SHIPMENT_ALLOCATE'
    GROUP BY `reference_id`
) alloc ON alloc.reference_id = s.id
SET
    s.`status` = 'DRAFT',
    s.`inventory_locked` = 0,
    s.`confirmed_at` = NULL,
    s.`confirmed_by` = NULL,
    s.`remark` = CASE
        WHEN s.`remark` IS NULL OR s.`remark` = '' THEN '[system repair 2026-03-09] confirmed shipment reset to draft because no shipment allocation movement exists'
        ELSE CONCAT(s.`remark`, '\n', '[system repair 2026-03-09] confirmed shipment reset to draft because no shipment allocation movement exists')
    END
WHERE s.`status` = 'CONFIRMED'
  AND s.`inventory_locked` = 1
  AND COALESCE(alloc.allocate_count, 0) = 0
  AND NOT EXISTS (
      SELECT 1
      FROM `inventory_movement` m2
      WHERE m2.`reference_type` = 'SHIPMENT'
        AND m2.`reference_id` = s.`id`
        AND m2.`movement_type` IN ('SHIPMENT_SHIP', 'PLATFORM_RECEIVE')
  );
