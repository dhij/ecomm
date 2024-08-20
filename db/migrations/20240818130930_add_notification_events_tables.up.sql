CREATE TABLE `notification_states` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `order_id` int NOT NULL,
  `state` enum('not sent', 'sent', 'failed') NOT NULL,
  `message` varchar(512),  
  `requested_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `completed_at` datetime
);

CREATE TABLE `notification_events_queue` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `user_email` varchar(256) NOT NULL,
  `order_status` varchar(256) NOT NULL,
  `order_id` int NOT NULL, 
  `state_id` int NOT NULL,
  `attempts` int,  
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime
);

ALTER TABLE `notification_states`
    ADD CONSTRAINT `notification_states_order_id_fk` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`);

ALTER TABLE `notification_events_queue`
    ADD CONSTRAINT `notification_events_queue_order_id_fk` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`),
    ADD CONSTRAINT `notification_events_queue_state_id_fk` FOREIGN KEY (`state_id`) REFERENCES `notification_states` (`id`);
