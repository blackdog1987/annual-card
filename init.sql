-- MySQL dump 10.13  Distrib 5.7.15, for linux-glibc2.5 (x86_64)
--
-- Host: localhost    Database: annual
-- ------------------------------------------------------
-- Server version	5.7.15-log

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `annual_card`
--

DROP TABLE IF EXISTS `annual_card`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `annual_card` (
  `card_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '年卡ID',
  `plan_id` int(10) DEFAULT NULL COMMENT '年卡计划ID',
  `card_name` varchar(255) DEFAULT '' COMMENT '年卡名称',
  `card_no` varchar(20) NOT NULL COMMENT '年卡卡号',
  `relation_uid` int(10) DEFAULT NULL COMMENT '关联用户ID',
  `bind_headimg` text COMMENT '绑定用户头像',
  `bind_name` varchar(32) DEFAULT NULL COMMENT '绑定用户姓名',
  `bind_contact` varchar(12) DEFAULT NULL COMMENT '绑定用户电话',
  `bind_idcard` char(18) DEFAULT NULL COMMENT '绑定身份证号',
  `card_passwd` char(10) DEFAULT NULL COMMENT '年卡密码',
  `expired_start` char(10) DEFAULT NULL COMMENT '有效期开始时间',
  `expired_stop` char(10) DEFAULT NULL COMMENT '有效期终止时间',
  `is_active` tinyint(1) DEFAULT '0' COMMENT '激活状态: 0: 未激活; 1: 激活',
  `created` int(10) NOT NULL COMMENT '年卡创建时间',
  `updated` int(10) NOT NULL COMMENT '年卡更新时间',
  PRIMARY KEY (`card_id`),
  UNIQUE KEY `card_no` (`card_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `annual_card`
--

LOCK TABLES `annual_card` WRITE;
/*!40000 ALTER TABLE `annual_card` DISABLE KEYS */;
/*!40000 ALTER TABLE `annual_card` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `annual_card_plan`
--

DROP TABLE IF EXISTS `annual_card_plan`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `annual_card_plan` (
  `cp_id` int(10) NOT NULL COMMENT '年卡计划ID',
  `channel` varchar(255) DEFAULT NULL COMMENT '渠道名称',
  `expired_start` int(10) DEFAULT NULL COMMENT '有效期开始时间',
  `expired_stop` int(10) DEFAULT NULL COMMENT '有效期截止时间',
  `create_num` int(10) DEFAULT NULL COMMENT '发放数量',
  `active_num` int(10) DEFAULT NULL COMMENT '激活数量',
  `created` int(10) DEFAULT NULL COMMENT '创建时间',
  `updated` int(10) DEFAULT NULL COMMENT '更新时间',
  `is_disabled` tinyint(1) DEFAULT '0' COMMENT '计划状态:\n0: 未关闭\n1: 已关闭',
  PRIMARY KEY (`cp_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='年卡计划';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `annual_card_plan`
--

LOCK TABLES `annual_card_plan` WRITE;
/*!40000 ALTER TABLE `annual_card_plan` DISABLE KEYS */;
/*!40000 ALTER TABLE `annual_card_plan` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `annual_card_usage_log`
--

DROP TABLE IF EXISTS `annual_card_usage_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `annual_card_usage_log` (
  `usage_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '使用记录ID',
  `card_id` int(10) unsigned NOT NULL COMMENT '年卡ID',
  `mch_id` int(10) unsigned NOT NULL COMMENT '商户ID',
  `usage_time` int(10) NOT NULL,
  PRIMARY KEY (`usage_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `annual_card_usage_log`
--

LOCK TABLES `annual_card_usage_log` WRITE;
/*!40000 ALTER TABLE `annual_card_usage_log` DISABLE KEYS */;
/*!40000 ALTER TABLE `annual_card_usage_log` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `banner`
--

DROP TABLE IF EXISTS `banner`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `banner` (
  `banner_id` int(10) NOT NULL AUTO_INCREMENT COMMENT 'banner ID',
  `name` varchar(64) DEFAULT NULL COMMENT '标题',
  `link` text COMMENT '跳转地址',
  `ordid` int(10) DEFAULT '255' COMMENT '排序',
  `created` int(10) NOT NULL,
  PRIMARY KEY (`banner_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `banner`
--

LOCK TABLES `banner` WRITE;
/*!40000 ALTER TABLE `banner` DISABLE KEYS */;
/*!40000 ALTER TABLE `banner` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `config`
--

DROP TABLE IF EXISTS `config`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `config` (
  `key` char(10) NOT NULL COMMENT '配置键值',
  `value` text COMMENT '配置值',
  `name` varchar(64) DEFAULT NULL COMMENT '配置名',
  PRIMARY KEY (`key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `config`
--

LOCK TABLES `config` WRITE;
/*!40000 ALTER TABLE `config` DISABLE KEYS */;
/*!40000 ALTER TABLE `config` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `manager`
--

DROP TABLE IF EXISTS `manager`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `manager` (
  `manager_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '管理人员ID',
  `name` varchar(64) NOT NULL COMMENT '姓名',
  `phone` char(11) NOT NULL COMMENT '手机号',
  `email` varchar(255) DEFAULT NULL COMMENT '邮箱',
  `passwd` char(40) NOT NULL COMMENT '密码',
  `created` int(10) NOT NULL,
  `updated` int(10) NOT NULL,
  `is_disabled` tinyint(1) DEFAULT '0' COMMENT '是否禁用 0: 否 1: 是',
  PRIMARY KEY (`manager_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `manager`
--

LOCK TABLES `manager` WRITE;
/*!40000 ALTER TABLE `manager` DISABLE KEYS */;
INSERT INTO `manager` VALUES (1,'超级管理员','15000000001',NULL,'f59bd65f7edafb087a81d4dca06c4910',1475424357,1475424357,0);
/*!40000 ALTER TABLE `manager` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `merchant`
--

DROP TABLE IF EXISTS `merchant`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `merchant` (
  `mch_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '商户ID',
  `mch_name` varchar(255) NOT NULL COMMENT '商户名称',
  `value` varchar(10) DEFAULT NULL COMMENT '价值',
  `consume` text COMMENT '消费特惠',
  `usage` text COMMENT '使用说明',
  `contact` varchar(255) DEFAULT NULL COMMENT '联系方式',
  `address` text COMMENT '联系地址',
  `introduce` text COMMENT '详细介绍',
  `cover` text COMMENT '封面图',
  `imgs` text COMMENT '商家介绍图',
  `created` int(10) NOT NULL,
  `updated` int(10) NOT NULL,
  `state` tinyint(1) DEFAULT '1' COMMENT '状态',
  PRIMARY KEY (`mch_id`),
  UNIQUE KEY `mch_name` (`mch_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `merchant`
--

LOCK TABLES `merchant` WRITE;
/*!40000 ALTER TABLE `merchant` DISABLE KEYS */;
/*!40000 ALTER TABLE `merchant` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `order_log`
--

DROP TABLE IF EXISTS `order_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `order_log` (
  `order_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '订单记录ID',
  `order_no` char(20) NOT NULL COMMENT '订单号',
  `transaction_id` varchar(64) DEFAULT NULL COMMENT '支付平台交易流水号',
  `total` int(10) unsigned DEFAULT '0' COMMENT '交易金额',
  `points` int(10) unsigned DEFAULT '0' COMMENT '积分抵扣',
  `is_pay` tinyint(1) DEFAULT '0' COMMENT '是否支付 0:等待支付 1: 支付完成',
  `created` int(10) NOT NULL,
  `updated` int(10) NOT NULL,
  PRIMARY KEY (`order_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `order_log`
--

LOCK TABLES `order_log` WRITE;
/*!40000 ALTER TABLE `order_log` DISABLE KEYS */;
/*!40000 ALTER TABLE `order_log` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `points_earning_log`
--

DROP TABLE IF EXISTS `points_earning_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `points_earning_log` (
  `earning_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '积分收入记录ID',
  `uid` int(10) unsigned NOT NULL COMMENT '用户ID',
  `relation_id` int(10) unsigned NOT NULL COMMENT '关联ID',
  `total` int(10) unsigned DEFAULT NULL COMMENT '积分个数',
  `type` tinyint(1) DEFAULT '0' COMMENT '收入类型 0: 推广;1: 交易',
  `created` int(10) NOT NULL,
  PRIMARY KEY (`earning_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `points_earning_log`
--

LOCK TABLES `points_earning_log` WRITE;
/*!40000 ALTER TABLE `points_earning_log` DISABLE KEYS */;
/*!40000 ALTER TABLE `points_earning_log` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `points_expend_log`
--

DROP TABLE IF EXISTS `points_expend_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `points_expend_log` (
  `expend_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '积分消费记录ID',
  `uid` int(10) unsigned NOT NULL COMMENT '用户ID',
  `relation_id` int(10) unsigned NOT NULL COMMENT '关联ID',
  `total` int(10) unsigned DEFAULT NULL COMMENT '积分个数',
  `type` tinyint(1) DEFAULT '0' COMMENT '收入类型 0: 抵扣;1: 兑换',
  `created` int(10) NOT NULL,
  PRIMARY KEY (`expend_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `points_expend_log`
--

LOCK TABLES `points_expend_log` WRITE;
/*!40000 ALTER TABLE `points_expend_log` DISABLE KEYS */;
/*!40000 ALTER TABLE `points_expend_log` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `points_log`
--

DROP TABLE IF EXISTS `points_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `points_log` (
  `log_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '记录ID',
  `uid` int(10) unsigned NOT NULL COMMENT '用户ID',
  `relation_uid` char(10) NOT NULL COMMENT '关联方',
  `relation_log_id` int(10) NOT NULL COMMENT '积分记录详情ID',
  `friendly_intro` text NOT NULL COMMENT '友好化描述',
  `total` int(10) DEFAULT '0' COMMENT '积分数量',
  `type` tinyint(1) NOT NULL COMMENT '积分交易类型 0: 支出; 1: 收入',
  `created` int(10) NOT NULL COMMENT '积分交易产生时间',
  PRIMARY KEY (`log_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `points_log`
--

LOCK TABLES `points_log` WRITE;
/*!40000 ALTER TABLE `points_log` DISABLE KEYS */;
/*!40000 ALTER TABLE `points_log` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `spread_logs`
--

DROP TABLE IF EXISTS `spread_logs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `spread_logs` (
  `log_id` int(10) NOT NULL AUTO_INCREMENT COMMENT '记录ID',
  `plan_id` int(10) DEFAULT NULL COMMENT '推广计划ID',
  `uid` int(10) DEFAULT NULL COMMENT '用户ID',
  `relation_id` int(10) DEFAULT NULL COMMENT '关联记录ID',
  `category` tinyint(1) DEFAULT NULL COMMENT '推广类型:\n1:注册\n2:销售',
  `commission` int(10) DEFAULT NULL COMMENT '佣金',
  `order_total` int(10) DEFAULT 0 COMMENT '订单金额',
  `created` int(10) DEFAULT NULL COMMENT '产生时间',
  PRIMARY KEY (`log_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='推广记录表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `spread_logs`
--

LOCK TABLES `spread_logs` WRITE;
/*!40000 ALTER TABLE `spread_logs` DISABLE KEYS */;
/*!40000 ALTER TABLE `spread_logs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `spread_plan`
--

DROP TABLE IF EXISTS `spread_plan`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `spread_plan` (
  `sp_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '推广计划ID',
  `name` varchar(255) NOT NULL COMMENT '推广计划名称',
  `channel` varchar(64) NOT NULL COMMENT '推广渠道',
  `contact` varchar(32) NOT NULL COMMENT '联系电话',
  `reg_commission` int(10) DEFAULT NULL,
  `sale_commission` int(10) DEFAULT NULL,
  `created` int(10) DEFAULT NULL,
  `updated` int(10) DEFAULT NULL COMMENT '停止时间',
  `is_disabled` tinyint(1) DEFAULT '0' COMMENT '是否关闭',
  `qrcode` text COMMENT '二维码',
  PRIMARY KEY (`sp_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COMMENT='推广计划';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `spread_plan`
--

LOCK TABLES `spread_plan` WRITE;
/*!40000 ALTER TABLE `spread_plan` DISABLE KEYS */;
INSERT INTO `spread_plan` VALUES (1,'测试推广计划','微信','023-11111111',10,1,1475856917,1475861175,0,'');
/*!40000 ALTER TABLE `spread_plan` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user`
--
CREATE TABLE `coupon_logs` (
  `log_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '记录ID',
  `uid` int(10) COMMENT '用户ID',
  `offset_price` int(10) DEFAULT 0 COMMENT '抵扣金额',
  `is_usage` tinyint(1) DEFAULT 0 COMMENT '是否使用',
  `created` int(10) DEFAULT 0,
  `updated` int(10) DEFAULT 0,
  PRIMARY KEY(`log_id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8;
DROP TABLE IF EXISTS `user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user` (
  `uid` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `spread_uid` int(10) DEFAULT NULL COMMENT '推广人',
  `nickname` varchar(32) DEFAULT '' COMMENT '昵称',
  `realname` varchar(32) DEFAULT '' COMMENT '真实姓名',
  `phone` varchar(12) DEFAULT '' COMMENT '手机号',
  `headimgurl` text COMMENT '头像',
  `sex` tinyint(1) DEFAULT NULL COMMENT '性别',
  `country` varchar(32) DEFAULT NULL COMMENT '国家',
  `province` varchar(32) DEFAULT NULL COMMENT '省份',
  `city` varchar(32) DEFAULT NULL COMMENT '城市',
  `area` varchar(32) DEFAULT NULL COMMENT '区县',
  `address` varchar(255) DEFAULT NULL COMMENT '地址',
  `balance` int(10) DEFAULT NULL COMMENT '账户余额',
  `points_balance` int(10) DEFAULT NULL COMMENT '积分余额',
  `points_earning` int(10) DEFAULT NULL COMMENT '积分收入合计',
  `pints_expend` int(10) DEFAULT NULL COMMENT '积分消费合集',
  `wx_openid` varchar(64) DEFAULT NULL COMMENT '微信开放平台id',
  `wx_unionid` varchar(64) DEFAULT NULL COMMENT '微信唯一ID',
  `created` int(10) NOT NULL,
  `updated` int(10) NOT NULL,
  PRIMARY KEY (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;
--
-- Dumping data for table `user`
--

LOCK TABLES `user` WRITE;
/*!40000 ALTER TABLE `user` DISABLE KEYS */;
/*!40000 ALTER TABLE `user` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2016-10-08  9:40:27
