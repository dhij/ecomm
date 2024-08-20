ALTER TABLE `orders`
	ADD COLUMN `status` ENUM('pending', 'shipped', 'delivered') NOT NULL DEFAULT 'pending';