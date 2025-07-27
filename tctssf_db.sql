-- phpMyAdmin SQL Dump
-- version 5.2.1
-- https://www.phpmyadmin.net/
--
-- Host: localhost
-- Generation Time: Jul 27, 2025 at 04:28 PM
-- Server version: 10.4.28-MariaDB
-- PHP Version: 8.2.4

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `tctssf`
--

-- --------------------------------------------------------

--
-- Table structure for table `loans`
--

CREATE TABLE `loans` (
  `id` int(11) NOT NULL,
  `user_id` int(11) NOT NULL,
  `amount` decimal(10,2) NOT NULL,
  `interest_rate` decimal(5,2) DEFAULT 5.00,
  `repayment_period` int(11) NOT NULL,
  `monthly_payment` decimal(10,2) DEFAULT 0.00,
  `remaining_balance` decimal(10,2) DEFAULT 0.00,
  `status` enum('pending','treasurer_approved','vice_president_approved','president_approved','approved','rejected','disbursed','completed') DEFAULT 'pending',
  `approved_by` int(11) DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `disbursed_at` timestamp NULL DEFAULT NULL,
  `treasurer_approved_by` int(11) DEFAULT NULL,
  `treasurer_approved_at` timestamp NULL DEFAULT NULL,
  `vice_president_approved_by` int(11) DEFAULT NULL,
  `vice_president_approved_at` timestamp NULL DEFAULT NULL,
  `president_approved_by` int(11) DEFAULT NULL,
  `president_approved_at` timestamp NULL DEFAULT NULL,
  `rejected_by` int(11) DEFAULT NULL,
  `rejected_at` timestamp NULL DEFAULT NULL,
  `rejection_reason` text DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `loans`
--

INSERT INTO `loans` (`id`, `user_id`, `amount`, `interest_rate`, `repayment_period`, `monthly_payment`, `remaining_balance`, `status`, `approved_by`, `created_at`, `disbursed_at`, `treasurer_approved_by`, `treasurer_approved_at`, `vice_president_approved_by`, `vice_president_approved_at`, `president_approved_by`, `president_approved_at`, `rejected_by`, `rejected_at`, `rejection_reason`) VALUES
(1, 7, 10000.00, 5.00, 12, 856.07, 10000.00, 'pending', NULL, '2025-07-03 14:47:37', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(2, 10, 40000.00, 5.00, 12, 3424.30, 40000.00, 'treasurer_approved', NULL, '2025-07-10 14:36:46', NULL, 3, '2025-07-10 14:40:26', NULL, NULL, NULL, NULL, NULL, NULL, NULL);

-- --------------------------------------------------------

--
-- Stand-in structure for view `loan_details_view`
-- (See below for the actual view)
--
CREATE TABLE `loan_details_view` (
`id` int(11)
,`user_id` int(11)
,`amount` decimal(10,2)
,`interest_rate` decimal(5,2)
,`repayment_period` int(11)
,`monthly_payment` decimal(10,2)
,`remaining_balance` decimal(10,2)
,`status` enum('pending','treasurer_approved','vice_president_approved','president_approved','approved','rejected','disbursed','completed')
,`approved_by` int(11)
,`created_at` timestamp
,`disbursed_at` timestamp
,`treasurer_approved_by` int(11)
,`treasurer_approved_at` timestamp
,`vice_president_approved_by` int(11)
,`vice_president_approved_at` timestamp
,`president_approved_by` int(11)
,`president_approved_at` timestamp
,`rejected_by` int(11)
,`rejected_at` timestamp
,`rejection_reason` text
,`member_first_name` varchar(100)
,`member_last_name` varchar(100)
,`account_number` varchar(20)
,`treasurer_name` varchar(100)
,`treasurer_last_name` varchar(100)
,`vp_name` varchar(100)
,`vp_last_name` varchar(100)
,`president_name` varchar(100)
,`president_last_name` varchar(100)
,`rejected_by_name` varchar(100)
,`rejected_by_last_name` varchar(100)
);

-- --------------------------------------------------------

--
-- Table structure for table `loan_repayments`
--

CREATE TABLE `loan_repayments` (
  `id` int(11) NOT NULL,
  `loan_id` int(11) NOT NULL,
  `amount` decimal(10,2) NOT NULL,
  `due_date` date NOT NULL,
  `paid_date` timestamp NULL DEFAULT NULL,
  `status` enum('pending','paid','overdue') DEFAULT 'pending',
  `created_at` timestamp NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table `salary_deduction_items`
--

CREATE TABLE `salary_deduction_items` (
  `id` int(11) NOT NULL,
  `list_id` int(11) NOT NULL,
  `user_id` int(11) NOT NULL,
  `account_number` varchar(20) NOT NULL,
  `member_name` varchar(200) NOT NULL,
  `monthly_commitment` decimal(10,2) DEFAULT 0.00,
  `social_contribution` decimal(10,2) DEFAULT 0.00,
  `loan_repayment` decimal(10,2) DEFAULT 0.00,
  `total_deduction` decimal(10,2) DEFAULT 0.00,
  `status` enum('pending','processed') DEFAULT 'pending'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `salary_deduction_items`
--

INSERT INTO `salary_deduction_items` (`id`, `list_id`, `user_id`, `account_number`, `member_name`, `monthly_commitment`, `social_contribution`, `loan_repayment`, `total_deduction`, `status`) VALUES
(5, 7, 6, '10051347', 'IRADUKUNDA Henriette Marie', 5000.00, 5000.00, 0.00, 10000.00, 'processed'),
(6, 7, 7, '10040901', 'MANZI Ivan Bright', 5000.00, 5000.00, 0.00, 10000.00, 'processed'),
(7, 7, 4, '10021881', 'Mugisha Emmanuel', 5000.00, 5000.00, 0.00, 10000.00, 'processed'),
(8, 7, 5, '10015566', 'Mukiza Innocent', 5000.00, 5000.00, 0.00, 10000.00, 'processed'),
(9, 8, 6, '10051347', 'IRADUKUNDA Henriette Marie', 5000.00, 5000.00, 0.00, 10000.00, 'processed'),
(10, 8, 7, '10040901', 'MANZI Ivan Bright', 5000.00, 5000.00, 0.00, 10000.00, 'pending'),
(11, 8, 4, '10021881', 'Mugisha Emmanuel', 5000.00, 5000.00, 0.00, 10000.00, 'processed'),
(12, 8, 5, '10015566', 'Mukiza Innocent', 5000.00, 5000.00, 0.00, 10000.00, 'processed'),
(13, 8, 10, '10059611', 'Niyonshuti Yves', 10000.00, 5000.00, 0.00, 15000.00, 'processed');

-- --------------------------------------------------------

--
-- Table structure for table `salary_deduction_lists`
--

CREATE TABLE `salary_deduction_lists` (
  `id` int(11) NOT NULL,
  `month_year` varchar(7) NOT NULL,
  `generated_by` int(11) NOT NULL,
  `status` enum('generated','sent_to_hr','processed') DEFAULT 'generated',
  `total_members` int(11) DEFAULT 0,
  `total_amount` decimal(12,2) DEFAULT 0.00,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `processed_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `salary_deduction_lists`
--

INSERT INTO `salary_deduction_lists` (`id`, `month_year`, `generated_by`, `status`, `total_members`, `total_amount`, `created_at`, `processed_at`) VALUES
(7, '2025-06', 3, 'processed', 4, 40000.00, '2025-06-28 23:05:29', '2025-06-28 23:16:53'),
(8, '2025-07', 3, 'processed', 5, 55000.00, '2025-07-10 14:25:29', '2025-07-10 14:31:12');

-- --------------------------------------------------------

--
-- Table structure for table `savings_accounts`
--

CREATE TABLE `savings_accounts` (
  `id` int(11) NOT NULL,
  `user_id` int(11) NOT NULL,
  `monthly_commitment` decimal(10,2) DEFAULT 0.00,
  `current_balance` decimal(10,2) DEFAULT 0.00,
  `social_contributions` decimal(10,2) DEFAULT 0.00,
  `last_contribution` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `savings_accounts`
--

INSERT INTO `savings_accounts` (`id`, `user_id`, `monthly_commitment`, `current_balance`, `social_contributions`, `last_contribution`, `created_at`) VALUES
(1, 4, 5000.00, 10000.00, 10000.00, '2025-07-10 14:31:12', '2025-06-28 11:08:21'),
(2, 5, 5000.00, 10000.00, 10000.00, '2025-07-10 14:31:12', '2025-06-28 11:08:21'),
(3, 6, 5000.00, 10000.00, 10000.00, '2025-07-10 14:31:12', '2025-06-28 11:08:21'),
(4, 7, 5000.00, 5000.00, 5000.00, '2025-06-28 23:16:53', '2025-06-28 11:08:21'),
(8, 10, 10000.00, 10000.00, 5000.00, '2025-07-10 14:31:12', '2025-07-10 14:18:52');

-- --------------------------------------------------------

--
-- Table structure for table `transactions`
--

CREATE TABLE `transactions` (
  `id` int(11) NOT NULL,
  `user_id` int(11) NOT NULL,
  `type` enum('savings','social_contribution','loan_disbursement','loan_repayment','salary_deduction','commitment_change') NOT NULL,
  `amount` decimal(10,2) NOT NULL,
  `description` text DEFAULT NULL,
  `reference_id` int(11) DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `transactions`
--

INSERT INTO `transactions` (`id`, `user_id`, `type`, `amount`, `description`, `reference_id`, `created_at`) VALUES
(1, 5, 'commitment_change', 20000.00, 'Monthly commitment updated to 20000 RWF, effective from 2025-07-01', NULL, '2025-06-28 10:53:15'),
(4, 6, 'savings', 5000.00, 'Monthly savings via salary deduction', 6, '2025-06-28 23:16:53'),
(5, 6, 'social_contribution', 5000.00, 'Social contribution via salary deduction', 6, '2025-06-28 23:16:53'),
(6, 7, 'savings', 5000.00, 'Monthly savings via salary deduction', 7, '2025-06-28 23:16:53'),
(7, 7, 'social_contribution', 5000.00, 'Social contribution via salary deduction', 7, '2025-06-28 23:16:53'),
(8, 4, 'savings', 5000.00, 'Monthly savings via salary deduction', 4, '2025-06-28 23:16:53'),
(9, 4, 'social_contribution', 5000.00, 'Social contribution via salary deduction', 4, '2025-06-28 23:16:53'),
(10, 5, 'savings', 5000.00, 'Monthly savings via salary deduction', 5, '2025-06-28 23:16:53'),
(11, 5, 'social_contribution', 5000.00, 'Social contribution via salary deduction', 5, '2025-06-28 23:16:53'),
(12, 10, 'commitment_change', 10000.00, 'Monthly commitment updated to 10000 RWF, effective from 2025-08-01', NULL, '2025-07-10 14:22:03'),
(13, 6, 'savings', 5000.00, 'Monthly savings via salary deduction', 6, '2025-07-10 14:31:12'),
(14, 6, 'social_contribution', 5000.00, 'Social contribution via salary deduction', 6, '2025-07-10 14:31:12'),
(15, 4, 'savings', 5000.00, 'Monthly savings via salary deduction', 4, '2025-07-10 14:31:12'),
(16, 4, 'social_contribution', 5000.00, 'Social contribution via salary deduction', 4, '2025-07-10 14:31:12'),
(17, 5, 'savings', 5000.00, 'Monthly savings via salary deduction', 5, '2025-07-10 14:31:12'),
(18, 5, 'social_contribution', 5000.00, 'Social contribution via salary deduction', 5, '2025-07-10 14:31:12'),
(19, 10, 'savings', 10000.00, 'Monthly savings via salary deduction', 10, '2025-07-10 14:31:12'),
(20, 10, 'social_contribution', 5000.00, 'Social contribution via salary deduction', 10, '2025-07-10 14:31:12');

-- --------------------------------------------------------

--
-- Table structure for table `users`
--

CREATE TABLE `users` (
  `id` int(11) NOT NULL,
  `account_number` varchar(20) DEFAULT NULL,
  `first_name` varchar(100) NOT NULL,
  `last_name` varchar(100) NOT NULL,
  `email` varchar(255) NOT NULL,
  `phone` varchar(20) DEFAULT NULL,
  `password_hash` varchar(255) NOT NULL,
  `role` enum('member','admin','superadmin','treasurer') DEFAULT 'member',
  `is_active` tinyint(1) DEFAULT 1,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `specific_role` varchar(50) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `users`
--

INSERT INTO `users` (`id`, `account_number`, `first_name`, `last_name`, `email`, `phone`, `password_hash`, `role`, `is_active`, `created_at`, `specific_role`) VALUES
(1, '10000000', 'Super', 'Admin', 'superadmin@tctssf.rw', NULL, '$2a$10$JvQUBnR6JTKCp/c089QIB.F/007dIowOtHsnsCZQ51bxDjfHQiQyK', 'superadmin', 1, '2025-06-27 11:20:18', 'super_admin'),
(2, '10000002', 'Default', 'Admin', 'admin@tctssf.rw', NULL, '$2a$10$8t9l7Q79sMw6BnlauIXgfuHxIfpn7qDWDsJz7UEAGUCq1fxPhl91C', 'admin', 1, '2025-06-27 11:20:18', 'general_admin'),
(3, '10000001', 'Default', 'Treasurer', 'treasurer@tctssf.rw', NULL, '$2a$10$ZviL6pvBQHIrWBD9wDq7HeBibsuZX7v9QO/bNlL7WWtih5MhRoRTO', 'treasurer', 1, '2025-06-27 11:20:18', 'treasurer'),
(4, '10021881', 'Mugisha', 'Emmanuel', 'mugisha@tctssf.rw', '0788212121', '$2a$10$K62ggYLDQK4R2ZvdvpYQxed.QPMNXPtMnX/xA3tc5utcRFPpIWi.S', 'member', 1, '2025-06-26 12:25:00', NULL),
(5, '10015566', 'Mukiza', 'Innocent', 'mukiza@tctssf.rw', '0788707070', '$2a$10$q8vk16aIT6iyzHRf9BKDfuQ8hN2H7ofimGcYhUXwCC4mb.4aBw..q', 'member', 1, '2025-06-26 12:25:54', NULL),
(6, '10051347', 'IRADUKUNDA', 'Henriette Marie', 'henriette@tctssf.rw', '0788123456', '$2a$10$r2Rff66tAluwl/BNcrIhkeA1qtMHsI3DWq.iCNTd1KEOz20Wn6bH6', 'member', 1, '2025-06-27 06:37:45', NULL),
(7, '10040901', 'MANZI', 'Ivan Bright', 'bright@tctssf.rw', '0788654321', '$2a$10$oCPMVA6iqqJ0rg.wXlGTzeJE3Y8iczI2s1xbOeKUqAKhEyBmXKY2i', 'member', 1, '2025-06-27 06:37:45', NULL),
(8, '10000003', 'Vice', 'President', 'vicepresident@tctssf.rw', NULL, '$2a$10$8t9l7Q79sMw6BnlauIXgfuHxIfpn7qDWDsJz7UEAGUCq1fxPhl91C', 'admin', 1, '2025-07-03 11:57:28', 'vice_president'),
(9, '10000004', 'President', 'TCTSSF', 'president@tctssf.rw', NULL, '$2a$10$8t9l7Q79sMw6BnlauIXgfuHxIfpn7qDWDsJz7UEAGUCq1fxPhl91C', 'admin', 1, '2025-07-03 11:57:28', 'president'),
(10, '10059611', 'Niyonshuti', 'Yves', 'yves@tctssf.rw', '0788882288', '$2a$10$mRspNcBlaogEIp9jhfC.auDVy0.yg8eXp75jhwq/sWu99dsL/Nr1K', 'member', 1, '2025-07-10 14:18:52', NULL);

-- --------------------------------------------------------

--
-- Structure for view `loan_details_view`
--
DROP TABLE IF EXISTS `loan_details_view`;

CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`localhost` SQL SECURITY DEFINER VIEW `loan_details_view`  AS SELECT `l`.`id` AS `id`, `l`.`user_id` AS `user_id`, `l`.`amount` AS `amount`, `l`.`interest_rate` AS `interest_rate`, `l`.`repayment_period` AS `repayment_period`, `l`.`monthly_payment` AS `monthly_payment`, `l`.`remaining_balance` AS `remaining_balance`, `l`.`status` AS `status`, `l`.`approved_by` AS `approved_by`, `l`.`created_at` AS `created_at`, `l`.`disbursed_at` AS `disbursed_at`, `l`.`treasurer_approved_by` AS `treasurer_approved_by`, `l`.`treasurer_approved_at` AS `treasurer_approved_at`, `l`.`vice_president_approved_by` AS `vice_president_approved_by`, `l`.`vice_president_approved_at` AS `vice_president_approved_at`, `l`.`president_approved_by` AS `president_approved_by`, `l`.`president_approved_at` AS `president_approved_at`, `l`.`rejected_by` AS `rejected_by`, `l`.`rejected_at` AS `rejected_at`, `l`.`rejection_reason` AS `rejection_reason`, `u`.`first_name` AS `member_first_name`, `u`.`last_name` AS `member_last_name`, `u`.`account_number` AS `account_number`, `t1`.`first_name` AS `treasurer_name`, `t1`.`last_name` AS `treasurer_last_name`, `t2`.`first_name` AS `vp_name`, `t2`.`last_name` AS `vp_last_name`, `t3`.`first_name` AS `president_name`, `t3`.`last_name` AS `president_last_name`, `r`.`first_name` AS `rejected_by_name`, `r`.`last_name` AS `rejected_by_last_name` FROM (((((`loans` `l` join `users` `u` on(`l`.`user_id` = `u`.`id`)) left join `users` `t1` on(`l`.`treasurer_approved_by` = `t1`.`id`)) left join `users` `t2` on(`l`.`vice_president_approved_by` = `t2`.`id`)) left join `users` `t3` on(`l`.`president_approved_by` = `t3`.`id`)) left join `users` `r` on(`l`.`rejected_by` = `r`.`id`)) ;

--
-- Indexes for dumped tables
--

--
-- Indexes for table `loans`
--
ALTER TABLE `loans`
  ADD PRIMARY KEY (`id`),
  ADD KEY `user_id` (`user_id`),
  ADD KEY `approved_by` (`approved_by`),
  ADD KEY `idx_loans_status` (`status`),
  ADD KEY `idx_loans_treasurer_approved_by` (`treasurer_approved_by`),
  ADD KEY `idx_loans_vice_president_approved_by` (`vice_president_approved_by`),
  ADD KEY `idx_loans_president_approved_by` (`president_approved_by`),
  ADD KEY `idx_loans_rejected_by` (`rejected_by`);

--
-- Indexes for table `loan_repayments`
--
ALTER TABLE `loan_repayments`
  ADD PRIMARY KEY (`id`),
  ADD KEY `loan_id` (`loan_id`);

--
-- Indexes for table `salary_deduction_items`
--
ALTER TABLE `salary_deduction_items`
  ADD PRIMARY KEY (`id`),
  ADD KEY `list_id` (`list_id`),
  ADD KEY `user_id` (`user_id`);

--
-- Indexes for table `salary_deduction_lists`
--
ALTER TABLE `salary_deduction_lists`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `unique_month_year` (`month_year`),
  ADD KEY `generated_by` (`generated_by`);

--
-- Indexes for table `savings_accounts`
--
ALTER TABLE `savings_accounts`
  ADD PRIMARY KEY (`id`),
  ADD KEY `user_id` (`user_id`);

--
-- Indexes for table `transactions`
--
ALTER TABLE `transactions`
  ADD PRIMARY KEY (`id`),
  ADD KEY `user_id` (`user_id`);

--
-- Indexes for table `users`
--
ALTER TABLE `users`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `email` (`email`),
  ADD UNIQUE KEY `account_number` (`account_number`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `loans`
--
ALTER TABLE `loans`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=3;

--
-- AUTO_INCREMENT for table `loan_repayments`
--
ALTER TABLE `loan_repayments`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `salary_deduction_items`
--
ALTER TABLE `salary_deduction_items`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=14;

--
-- AUTO_INCREMENT for table `salary_deduction_lists`
--
ALTER TABLE `salary_deduction_lists`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=9;

--
-- AUTO_INCREMENT for table `savings_accounts`
--
ALTER TABLE `savings_accounts`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=9;

--
-- AUTO_INCREMENT for table `transactions`
--
ALTER TABLE `transactions`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=21;

--
-- AUTO_INCREMENT for table `users`
--
ALTER TABLE `users`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=11;

--
-- Constraints for dumped tables
--

--
-- Constraints for table `loans`
--
ALTER TABLE `loans`
  ADD CONSTRAINT `fk_loans_president_approved_by` FOREIGN KEY (`president_approved_by`) REFERENCES `users` (`id`) ON DELETE SET NULL,
  ADD CONSTRAINT `fk_loans_rejected_by` FOREIGN KEY (`rejected_by`) REFERENCES `users` (`id`) ON DELETE SET NULL,
  ADD CONSTRAINT `fk_loans_treasurer_approved_by` FOREIGN KEY (`treasurer_approved_by`) REFERENCES `users` (`id`) ON DELETE SET NULL,
  ADD CONSTRAINT `fk_loans_vice_president_approved_by` FOREIGN KEY (`vice_president_approved_by`) REFERENCES `users` (`id`) ON DELETE SET NULL,
  ADD CONSTRAINT `loans_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `loans_ibfk_2` FOREIGN KEY (`approved_by`) REFERENCES `users` (`id`) ON DELETE SET NULL;

--
-- Constraints for table `loan_repayments`
--
ALTER TABLE `loan_repayments`
  ADD CONSTRAINT `loan_repayments_ibfk_1` FOREIGN KEY (`loan_id`) REFERENCES `loans` (`id`) ON DELETE CASCADE;

--
-- Constraints for table `salary_deduction_items`
--
ALTER TABLE `salary_deduction_items`
  ADD CONSTRAINT `salary_deduction_items_ibfk_1` FOREIGN KEY (`list_id`) REFERENCES `salary_deduction_lists` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `salary_deduction_items_ibfk_2` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE;

--
-- Constraints for table `salary_deduction_lists`
--
ALTER TABLE `salary_deduction_lists`
  ADD CONSTRAINT `salary_deduction_lists_ibfk_1` FOREIGN KEY (`generated_by`) REFERENCES `users` (`id`);

--
-- Constraints for table `savings_accounts`
--
ALTER TABLE `savings_accounts`
  ADD CONSTRAINT `savings_accounts_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE;

--
-- Constraints for table `transactions`
--
ALTER TABLE `transactions`
  ADD CONSTRAINT `transactions_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
