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