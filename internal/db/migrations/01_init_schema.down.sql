-- Down Migration: Remove entire schema
-- This drops all tables and enums in reverse dependency order

-- 1. Drop all foreign key constraints first
ALTER TABLE "seller_profiles" DROP CONSTRAINT IF EXISTS "seller_profiles_seller_id_fkey";
ALTER TABLE "user_addresses" DROP CONSTRAINT IF EXISTS "user_addresses_user_id_fkey";
ALTER TABLE "gundams" DROP CONSTRAINT IF EXISTS "gundams_owner_id_fkey";
ALTER TABLE "gundams" DROP CONSTRAINT IF EXISTS "gundams_grade_id_fkey";
ALTER TABLE "gundam_accessories" DROP CONSTRAINT IF EXISTS "gundam_accessories_gundam_id_fkey";
ALTER TABLE "gundam_images" DROP CONSTRAINT IF EXISTS "gundam_images_gundam_id_fkey";
ALTER TABLE "carts" DROP CONSTRAINT IF EXISTS "carts_user_id_fkey";
ALTER TABLE "cart_items" DROP CONSTRAINT IF EXISTS "cart_items_cart_id_fkey";
ALTER TABLE "cart_items" DROP CONSTRAINT IF EXISTS "cart_items_gundam_id_fkey";
ALTER TABLE "seller_subscriptions" DROP CONSTRAINT IF EXISTS "seller_subscriptions_plan_id_fkey";
ALTER TABLE "seller_subscriptions" DROP CONSTRAINT IF EXISTS "seller_subscriptions_seller_id_fkey";
ALTER TABLE "orders" DROP CONSTRAINT IF EXISTS "orders_buyer_id_fkey";
ALTER TABLE "orders" DROP CONSTRAINT IF EXISTS "orders_seller_id_fkey";
ALTER TABLE "orders" DROP CONSTRAINT IF EXISTS "orders_canceled_by_fkey";
ALTER TABLE "order_items" DROP CONSTRAINT IF EXISTS "order_items_order_id_fkey";
ALTER TABLE "order_items" DROP CONSTRAINT IF EXISTS "order_items_gundam_id_fkey";
ALTER TABLE "delivery_information" DROP CONSTRAINT IF EXISTS "delivery_information_user_id_fkey";
ALTER TABLE "order_deliveries" DROP CONSTRAINT IF EXISTS "order_deliveries_from_delivery_id_fkey";
ALTER TABLE "order_deliveries" DROP CONSTRAINT IF EXISTS "order_deliveries_to_delivery_id_fkey";
ALTER TABLE "order_deliveries" DROP CONSTRAINT IF EXISTS "order_deliveries_order_id_fkey";
ALTER TABLE "wallets" DROP CONSTRAINT IF EXISTS "wallets_user_id_fkey";
ALTER TABLE "wallet_entries" DROP CONSTRAINT IF EXISTS "wallet_entries_wallet_id_fkey";
ALTER TABLE "order_transactions" DROP CONSTRAINT IF EXISTS "order_transactions_buyer_entry_id_fkey";
ALTER TABLE "order_transactions" DROP CONSTRAINT IF EXISTS "order_transactions_seller_entry_id_fkey";
ALTER TABLE "order_transactions" DROP CONSTRAINT IF EXISTS "order_transactions_order_id_fkey";
ALTER TABLE "payment_transactions" DROP CONSTRAINT IF EXISTS "payment_transactions_user_id_fkey";
ALTER TABLE "exchange_posts" DROP CONSTRAINT IF EXISTS "exchange_posts_user_id_fkey";
ALTER TABLE "exchange_post_items" DROP CONSTRAINT IF EXISTS "exchange_post_items_gundam_id_fkey";
ALTER TABLE "exchange_post_items" DROP CONSTRAINT IF EXISTS "exchange_post_items_post_id_fkey";
ALTER TABLE "exchange_offers" DROP CONSTRAINT IF EXISTS "exchange_offers_post_id_fkey";
ALTER TABLE "exchange_offers" DROP CONSTRAINT IF EXISTS "exchange_offers_offerer_id_fkey";
ALTER TABLE "exchange_offers" DROP CONSTRAINT IF EXISTS "exchange_offers_payer_id_fkey";
ALTER TABLE "exchange_offer_notes" DROP CONSTRAINT IF EXISTS "exchange_offer_notes_offer_id_fkey";
ALTER TABLE "exchange_offer_notes" DROP CONSTRAINT IF EXISTS "exchange_offer_notes_user_id_fkey";
ALTER TABLE "exchange_offer_items" DROP CONSTRAINT IF EXISTS "exchange_offer_items_gundam_id_fkey";
ALTER TABLE "exchange_offer_items" DROP CONSTRAINT IF EXISTS "exchange_offer_items_offer_id_fkey";
ALTER TABLE "exchanges" DROP CONSTRAINT IF EXISTS "exchanges_poster_from_delivery_id_fkey";
ALTER TABLE "exchanges" DROP CONSTRAINT IF EXISTS "exchanges_poster_to_delivery_id_fkey";
ALTER TABLE "exchanges" DROP CONSTRAINT IF EXISTS "exchanges_offerer_from_delivery_id_fkey";
ALTER TABLE "exchanges" DROP CONSTRAINT IF EXISTS "exchanges_offerer_to_delivery_id_fkey";
ALTER TABLE "exchanges" DROP CONSTRAINT IF EXISTS "exchanges_payer_id_fkey";
ALTER TABLE "exchanges" DROP CONSTRAINT IF EXISTS "exchanges_canceled_by_fkey";
ALTER TABLE "exchanges" DROP CONSTRAINT IF EXISTS "exchanges_poster_order_id_fkey";
ALTER TABLE "exchanges" DROP CONSTRAINT IF EXISTS "exchanges_offerer_order_id_fkey";
ALTER TABLE "exchange_items" DROP CONSTRAINT IF EXISTS "exchange_items_exchange_id_fkey";
ALTER TABLE "exchange_items" DROP CONSTRAINT IF EXISTS "exchange_items_gundam_id_fkey";
ALTER TABLE "exchange_items" DROP CONSTRAINT IF EXISTS "exchange_items_owner_id_fkey";
ALTER TABLE "auction_requests" DROP CONSTRAINT IF EXISTS "auction_requests_gundam_id_fkey";
ALTER TABLE "auction_requests" DROP CONSTRAINT IF EXISTS "auction_requests_seller_id_fkey";
ALTER TABLE "auction_requests" DROP CONSTRAINT IF EXISTS "auction_requests_rejected_by_fkey";
ALTER TABLE "auction_requests" DROP CONSTRAINT IF EXISTS "auction_requests_approved_by_fkey";
ALTER TABLE "auctions" DROP CONSTRAINT IF EXISTS "auctions_winning_bid_id_fkey";
ALTER TABLE "auctions" DROP CONSTRAINT IF EXISTS "auctions_request_id_fkey";
ALTER TABLE "auctions" DROP CONSTRAINT IF EXISTS "auctions_gundam_id_fkey";
ALTER TABLE "auctions" DROP CONSTRAINT IF EXISTS "auctions_seller_id_fkey";
ALTER TABLE "auctions" DROP CONSTRAINT IF EXISTS "auctions_order_id_fkey";
ALTER TABLE "auctions" DROP CONSTRAINT IF EXISTS "auctions_canceled_by_fkey";
ALTER TABLE "auction_bids" DROP CONSTRAINT IF EXISTS "auction_bids_auction_id_fkey";
ALTER TABLE "auction_bids" DROP CONSTRAINT IF EXISTS "auction_bids_bidder_id_fkey";
ALTER TABLE "auction_bids" DROP CONSTRAINT IF EXISTS "auction_bids_participant_id_fkey";
ALTER TABLE "auction_participants" DROP CONSTRAINT IF EXISTS "auction_participants_deposit_entry_id_fkey";
ALTER TABLE "auction_participants" DROP CONSTRAINT IF EXISTS "auction_participants_auction_id_fkey";
ALTER TABLE "auction_participants" DROP CONSTRAINT IF EXISTS "auction_participants_user_id_fkey";

-- 2. Drop all indexes
DROP INDEX IF EXISTS "user_addresses_user_id_is_primary_idx";
DROP INDEX IF EXISTS "user_addresses_user_id_is_pickup_address_idx";
DROP INDEX IF EXISTS "unique_cart_item";
DROP INDEX IF EXISTS "idx_seller_active_subscription";
DROP INDEX IF EXISTS "wallets_user_id_idx";
DROP INDEX IF EXISTS "wallet_entries_wallet_id_idx";
DROP INDEX IF EXISTS "wallet_entries_reference_id_reference_type_idx";
DROP INDEX IF EXISTS "order_transactions_order_id_idx";
DROP INDEX IF EXISTS "payment_transactions_provider_provider_transaction_id_idx";
DROP INDEX IF EXISTS "payment_transactions_user_id_status_idx";
DROP INDEX IF EXISTS "exchange_posts_user_id_idx";
DROP INDEX IF EXISTS "exchange_posts_status_idx";
DROP INDEX IF EXISTS "exchange_posts_created_at_idx";
DROP INDEX IF EXISTS "exchange_post_items_post_id_gundam_id_idx";
DROP INDEX IF EXISTS "exchange_offers_post_id_idx";
DROP INDEX IF EXISTS "exchange_offers_offerer_id_idx";
DROP INDEX IF EXISTS "exchange_offers_created_at_idx";
DROP INDEX IF EXISTS "unique_exchange_offer";
DROP INDEX IF EXISTS "exchange_offer_notes_offer_id_idx";
DROP INDEX IF EXISTS "exchange_offer_notes_user_id_idx";
DROP INDEX IF EXISTS "exchange_offer_notes_created_at_idx";
DROP INDEX IF EXISTS "exchange_offer_items_offer_id_gundam_id_idx";
DROP INDEX IF EXISTS "exchanges_poster_order_id_idx";
DROP INDEX IF EXISTS "exchanges_offerer_order_id_idx";
DROP INDEX IF EXISTS "exchanges_status_idx";
DROP INDEX IF EXISTS "exchange_items_exchange_id_gundam_id_idx";
DROP INDEX IF EXISTS "auction_requests_seller_id_status_idx";
DROP INDEX IF EXISTS "auction_requests_status_created_at_idx";
DROP INDEX IF EXISTS "auction_requests_gundam_id_status_idx";
DROP INDEX IF EXISTS "auctions_status_start_time_idx";
DROP INDEX IF EXISTS "auctions_status_end_time_idx";
DROP INDEX IF EXISTS "auctions_seller_id_status_idx";
DROP INDEX IF EXISTS "auctions_gundam_id_idx";
DROP INDEX IF EXISTS "auctions_current_price_idx";
DROP INDEX IF EXISTS "auction_bids_auction_id_bidder_id_idx";
DROP INDEX IF EXISTS "auction_bids_auction_id_amount_idx";
DROP INDEX IF EXISTS "auction_bids_participant_id_idx";
DROP INDEX IF EXISTS "auction_bids_bidder_id_created_at_idx";
DROP INDEX IF EXISTS "auction_bids_auction_id_created_at_idx";
DROP INDEX IF EXISTS "auction_participants_auction_id_user_id_idx";
DROP INDEX IF EXISTS "auction_participants_user_id_created_at_idx";

-- 3. Drop all tables in reverse dependency order
DROP TABLE IF EXISTS "auction_participants";
DROP TABLE IF EXISTS "auction_bids";
DROP TABLE IF EXISTS "auctions";
DROP TABLE IF EXISTS "auction_requests";
DROP TABLE IF EXISTS "exchange_items";
DROP TABLE IF EXISTS "exchanges";
DROP TABLE IF EXISTS "exchange_offer_items";
DROP TABLE IF EXISTS "exchange_offer_notes";
DROP TABLE IF EXISTS "exchange_offers";
DROP TABLE IF EXISTS "exchange_post_items";
DROP TABLE IF EXISTS "exchange_posts";
DROP TABLE IF EXISTS "payment_transactions";
DROP TABLE IF EXISTS "order_transactions";
DROP TABLE IF EXISTS "wallet_entries";
DROP TABLE IF EXISTS "wallets";
DROP TABLE IF EXISTS "order_deliveries";
DROP TABLE IF EXISTS "delivery_information";
DROP TABLE IF EXISTS "order_items";
DROP TABLE IF EXISTS "orders";
DROP TABLE IF EXISTS "seller_subscriptions";
DROP TABLE IF EXISTS "subscription_plans";
DROP TABLE IF EXISTS "cart_items";
DROP TABLE IF EXISTS "carts";
DROP TABLE IF EXISTS "gundam_images";
DROP TABLE IF EXISTS "gundam_accessories";
DROP TABLE IF EXISTS "gundam_grades";
DROP TABLE IF EXISTS "gundams";
DROP TABLE IF EXISTS "user_addresses";
DROP TABLE IF EXISTS "seller_profiles";
DROP TABLE IF EXISTS "users";

-- 4. Drop all enum types
DROP TYPE IF EXISTS "auction_status";
DROP TYPE IF EXISTS "auction_request_status";
DROP TYPE IF EXISTS "exchange_status";
DROP TYPE IF EXISTS "exchange_post_status";
DROP TYPE IF EXISTS "payment_transaction_type";
DROP TYPE IF EXISTS "payment_transaction_status";
DROP TYPE IF EXISTS "payment_transaction_provider";
DROP TYPE IF EXISTS "order_transaction_status";
DROP TYPE IF EXISTS "wallet_affected_field";
DROP TYPE IF EXISTS "wallet_entry_status";
DROP TYPE IF EXISTS "wallet_reference_type";
DROP TYPE IF EXISTS "wallet_entry_type";
DROP TYPE IF EXISTS "delivery_overral_status";
DROP TYPE IF EXISTS "order_type";
DROP TYPE IF EXISTS "payment_method";
DROP TYPE IF EXISTS "order_status";
DROP TYPE IF EXISTS "gundam_status";
DROP TYPE IF EXISTS "gundam_scale";
DROP TYPE IF EXISTS "gundam_condition";
DROP TYPE IF EXISTS "user_role";