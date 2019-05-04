-- File Table
CREATE TABLE `tbl_file` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `file_sha1` char(40) NOT NULL DEFAULT '' COMMENT 'File Hash',
  `file_name` varchar(256) NOT NULL DEFAULT '' COMMENT 'File Name',
  `file_size` bigint(20) DEFAULT '0' COMMENT 'File Size',
  `file_addr` varchar(1024) NOT NULL DEFAULT '' COMMENT 'File Address',
  `create_at` datetime default NOW() COMMENT 'Created Time',
  `update_at` datetime default NOW() on update current_timestamp() COMMENT 'Updated Time',
  `status` int(11) NOT NULL DEFAULT '0' COMMENT 'State(Enable/Disable/Deleted etc.)',
  `ext1` int(11) DEFAULT '0' COMMENT 'Alternate Field 1',
  `ext2` text COMMENT 'Alternate Field 2',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_file_hash` (`file_sha1`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- User Table
CREATE TABLE `tbl_user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_name` varchar(64) NOT NULL DEFAULT '' COMMENT 'User Name',
  `user_pwd` varchar(256) NOT NULL DEFAULT '' COMMENT 'User Password with salted hash',
  `email` varchar(64) DEFAULT '' COMMENT 'Email',
  `phone` varchar(128) DEFAULT '' COMMENT 'Phone Number',
  `email_validated` tinyint(1) DEFAULT 0 COMMENT 'Email Validation Status',
  `phone_validated` tinyint(1) DEFAULT 0 COMMENT 'Phone Validation Status',
  `signup_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'Register Date',
  `last_active` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Last Sign In Date',
  `profile` text COMMENT 'Profile Status',
  `status` int(11) NOT NULL DEFAULT '0' COMMENT 'User Account Status(Activate/Deactivate/Lock/Mark as Deleted/etc.)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_username` (`user_name`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4;

-- User Token Table
CREATE TABLE `tbl_user_token` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_name` varchar(64) NOT NULL DEFAULT '' COMMENT 'User Name',
  `user_token` char(40) NOT NULL DEFAULT '' COMMENT 'User Sign In Token',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_username` (`user_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- User File Table
CREATE TABLE `tbl_user_file` (
  `id` int(11) NOT NULL PRIMARY KEY AUTO_INCREMENT,
  `user_name` varchar(64) NOT NULL,
  `file_sha1` varchar(64) NOT NULL DEFAULT '' COMMENT 'File Hash',
  `file_size` bigint(20) DEFAULT '0' COMMENT 'File Size',
  `file_name` varchar(256) NOT NULL DEFAULT '' COMMENT 'File Name',
  `upload_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'Upload Time',
  `last_update` datetime DEFAULT CURRENT_TIMESTAMP 
          ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated Time',
  `status` int(11) NOT NULL DEFAULT '0' COMMENT 'State(0 is Enable, 1 is Deleted, 2 is Disable)',
  KEY `idx_status` (`status`),
  KEY `idx_user_id` (`user_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;