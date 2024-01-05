package migrations

// v0Initial bootstraps the schema with an initial migration.
var v0Initial = Migration{
	Name: "Initial",
	SQL: `
		SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';

		CREATE TABLE pricing_products (
			id INT(8) UNSIGNED AUTO_INCREMENT,
			provider VARCHAR(16) NOT NULL,
			sku VARCHAR(100) NOT NULL,
			meter_id VARCHAR(100) NOT NULL,	
			location VARCHAR(100) NOT NULL,
			service VARCHAR(100) NOT NULL,
			family VARCHAR(100) NULL,
			attributes JSON NOT NULL,
			PRIMARY KEY (id),
			UNIQUE KEY (provider, sku, meter_id),
			INDEX (provider, location, service, family)
		);

		CREATE TABLE pricing_product_prices (
			id INT(8) UNSIGNED AUTO_INCREMENT,
			product_id INT(8) UNSIGNED NOT NULL,
			hash VARCHAR(32) NOT NULL,
			currency VARCHAR(16) NOT NULL,
			unit VARCHAR(32) NOT NULL,
			price DECIMAL(24,10) NOT NULL,
			attributes JSON NOT NULL,
			PRIMARY KEY (id),
			UNIQUE KEY (product_id, hash),
			FOREIGN KEY (product_id) REFERENCES pricing_products (id)
		);

		CREATE TABLE ingestion_jobs (
			id INT(8) UNSIGNED AUTO_INCREMENT,
		    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			provider VARCHAR(16) NOT NULL,
			location VARCHAR(100) NOT NULL,
			service VARCHAR(100) NOT NULL,
			status VARCHAR(16) NOT NULL,
			error_msg VARCHAR(500),
			PRIMARY KEY (id)
		);
	`,
}
