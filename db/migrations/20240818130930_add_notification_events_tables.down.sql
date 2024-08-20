ALTER TABLE `notification_states`
    DROP FOREIGN KEY `notification_states_order_id_fk`,
    DROP COLUMN `order_id`;

ALTER TABLE `notification_events_queue`
    DROP FOREIGN KEY `notification_events_queue_order_id_fk`,
    DROP FOREIGN KEY `notification_events_queue_state_id_fk`,
    DROP COLUMN `order_id`,
    DROP COLUMN `state_id`;

DROP TABLE IF EXISTS notification_events_queue;
DROP TABLE IF EXISTS notification_states;