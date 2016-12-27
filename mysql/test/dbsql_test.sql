-- MySQL dump 10.13  Distrib 5.5.53, for debian-linux-gnu (x86_64)
--
-- Host: localhost    Database: dbsql_test
-- ------------------------------------------------------
-- Server version	5.5.53-0+deb8u1

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `abc`
--

DROP TABLE IF EXISTS `abc`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `abc` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` char(12) NOT NULL,
  `description` varchar(20) NOT NULL,
  `tiny` tinyint(4) DEFAULT '3',
  `small` smallint(6) DEFAULT '11',
  `medium` mediumint(9) DEFAULT '42',
  `ger` int(11) DEFAULT NULL,
  `big` bigint(20) DEFAULT NULL,
  `cost` decimal(10,0) DEFAULT NULL,
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `abc`
--

LOCK TABLES `abc` WRITE;
/*!40000 ALTER TABLE `abc` DISABLE KEYS */;
/*!40000 ALTER TABLE `abc` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `abc_nn`
--

DROP TABLE IF EXISTS `abc_nn`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `abc_nn` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` char(12) NOT NULL,
  `description` varchar(20) NOT NULL,
  `tiny` tinyint(4) NOT NULL,
  `small` smallint(6) NOT NULL,
  `medium` mediumint(9) NOT NULL,
  `ger` int(11) NOT NULL,
  `big` bigint(20) NOT NULL,
  `cost` decimal(10,0) NOT NULL,
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `abc_nn`
--

LOCK TABLES `abc_nn` WRITE;
/*!40000 ALTER TABLE `abc_nn` DISABLE KEYS */;
/*!40000 ALTER TABLE `abc_nn` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `def`
--

DROP TABLE IF EXISTS `def`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `def` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `d_date` date DEFAULT NULL,
  `d_datetime` datetime DEFAULT NULL,
  `d_time` time DEFAULT NULL,
  `d_year` year(4) DEFAULT NULL,
  `size` enum('small','medium','large') DEFAULT NULL,
  `a_set` set('a','b','c') DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `def`
--

LOCK TABLES `def` WRITE;
/*!40000 ALTER TABLE `def` DISABLE KEYS */;
/*!40000 ALTER TABLE `def` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `def_nn`
--

DROP TABLE IF EXISTS `def_nn`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `def_nn` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `d_date` date NOT NULL,
  `d_datetime` datetime NOT NULL,
  `d_time` time NOT NULL,
  `d_year` year(4) NOT NULL,
  `size` enum('small','med','large') NOT NULL,
  `a_set` set('1','2','3') NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `def_nn`
--

LOCK TABLES `def_nn` WRITE;
/*!40000 ALTER TABLE `def_nn` DISABLE KEYS */;
/*!40000 ALTER TABLE `def_nn` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ghi`
--

DROP TABLE IF EXISTS `ghi`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ghi` (
  `tiny_stuff` tinyblob,
  `stuff` blob,
  `med_stuff` mediumblob,
  `long_stuff` longblob
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ghi`
--

LOCK TABLES `ghi` WRITE;
/*!40000 ALTER TABLE `ghi` DISABLE KEYS */;
/*!40000 ALTER TABLE `ghi` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ghi_nn`
--

DROP TABLE IF EXISTS `ghi_nn`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ghi_nn` (
  `tiny_stuff` tinyblob NOT NULL,
  `stuff` blob NOT NULL,
  `med_stuff` mediumblob NOT NULL,
  `long_stuff` longblob NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ghi_nn`
--

LOCK TABLES `ghi_nn` WRITE;
/*!40000 ALTER TABLE `ghi_nn` DISABLE KEYS */;
/*!40000 ALTER TABLE `ghi_nn` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `jkl`
--

DROP TABLE IF EXISTS `jkl`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `jkl` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `tiny_txt` tinytext,
  `txt` text,
  `med_txt` mediumtext,
  `long_txt` longtext,
  `bin` binary(3) DEFAULT NULL,
  `var_bin` varbinary(12) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=ascii;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `jkl`
--

LOCK TABLES `jkl` WRITE;
/*!40000 ALTER TABLE `jkl` DISABLE KEYS */;
/*!40000 ALTER TABLE `jkl` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `jkl_nn`
--

DROP TABLE IF EXISTS `jkl_nn`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `jkl_nn` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `tiny_txt` tinytext NOT NULL,
  `txt` text NOT NULL,
  `med_txt` mediumtext NOT NULL,
  `long_txt` longtext NOT NULL,
  `bin` binary(3) NOT NULL,
  `var_bin` varbinary(12) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=ascii;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `jkl_nn`
--

LOCK TABLES `jkl_nn` WRITE;
/*!40000 ALTER TABLE `jkl_nn` DISABLE KEYS */;
/*!40000 ALTER TABLE `jkl_nn` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `mno`
--

DROP TABLE IF EXISTS `mno`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `mno` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `geo` geometry DEFAULT NULL,
  `pt` point DEFAULT NULL,
  `lstring` linestring DEFAULT NULL,
  `poly` polygon DEFAULT NULL,
  `multi_pt` multipoint DEFAULT NULL,
  `multi_lstring` multilinestring DEFAULT NULL,
  `multi_polygon` multipolygon DEFAULT NULL,
  `geo_collection` geometrycollection DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `mno`
--

LOCK TABLES `mno` WRITE;
/*!40000 ALTER TABLE `mno` DISABLE KEYS */;
/*!40000 ALTER TABLE `mno` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `mno_nn`
--

DROP TABLE IF EXISTS `mno_nn`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `mno_nn` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `geo` geometry NOT NULL,
  `pt` point NOT NULL,
  `lstring` linestring NOT NULL,
  `poly` polygon NOT NULL,
  `multi_pt` multipoint NOT NULL,
  `multi_lstring` multilinestring NOT NULL,
  `multi_polygon` multipolygon NOT NULL,
  `geo_collection` geometrycollection NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `mno_nn`
--

LOCK TABLES `mno_nn` WRITE;
/*!40000 ALTER TABLE `mno_nn` DISABLE KEYS */;
/*!40000 ALTER TABLE `mno_nn` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2016-12-27 16:54:07
