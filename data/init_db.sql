USE `csp_reports`;

# This will create the schema for the CSP Reports
# https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Security-Policy-Report-Only
CREATE TABLE `reports` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `document_uri` varchar(250) NOT NULL,
  `referrer` varchar(250) DEFAULT NULL,
  `violated_directive` varchar(250) DEFAULT NULL,
  `effective_directive` varchar(250) DEFAULT NULL,
  `original_policy` varchar(500) DEFAULT NULL,
  `disposition` varchar(25) DEFAULT NULL,
  `blocked_uri` varchar(250) DEFAULT NULL,
  `line_number` int(11) DEFAULT NULL,
  `column_number` int(11) DEFAULT NULL,
  `source_file` varchar(250) DEFAULT NULL,
  `status_code` int(4) DEFAULT 0,
  `script_sample` varchar(150) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;