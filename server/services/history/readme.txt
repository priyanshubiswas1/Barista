session service (rmq publish): SessionCreated event
    history service (rmq consume): onSessionCreated event -> INSERT INTO session_summary (
                                                    session_id,
                                                    session_number,
                                                    session_create_date,
                                                    session_update_date
                                                )
                                                VALUES (
                                                    '<session_id>',
                                                    <session_number>,
                                                    now(),
                                                    now()
                                                );

-------------------------------------------------------------------------------------------------------------------------------

session service (rmq publish): SessionClosed event
    history service (rmq consume): onSessionClosed event -> UPDATE session_summary
                                                                SET session_update_date = now(),
                                                                    updated_at = now()
                                                                WHERE session_id = '<session_id>';

-------------------------------------------------------------------------------------------------------------------------------

order service (rmq publish): OrderItemAdded event
    history service (rmq consume): onOrderItemAdded event -> UPDATE session_summary
                                                                SET 
                                                                    ordered_items = ordered_items || jsonb_build_object(
                                                                        'orderId', '<order_id>',
                                                                        'order_itemid', '<order_item_id>'
                                                                    ),
                                                                    session_total_amount = (
                                                                        session_total_amount + <item_price> * <quantity>
                                                                    ),
                                                                    session_update_date = now(),
                                                                    updated_at = now()
                                                                WHERE session_id = '<session_id>';

-------------------------------------------------------------------------------------------------------------------------------

order service (rmq publish): OrderItemRemoved event
    history service (rmq consume): onOrderItemRemoved event -> UPDATE session_summary
                                                                SET 
                                                                    ordered_items = (
                                                                        SELECT jsonb_agg(elem)
                                                                        FROM jsonb_array_elements(ordered_items) AS elem
                                                                        WHERE elem->>'order_itemid' != '<order_item_id>'
                                                                    ),
                                                                    session_total_amount = (
                                                                        session_total_amount - <item_price> * <quantity>
                                                                    ),
                                                                    session_update_date = now(),
                                                                    updated_at = now()
                                                                WHERE session_id = '<session_id>';

-------------------------------------------------------------------------------------------------------------------------------

SELECT *
FROM session_summary
WHERE ordered_items @> '[{"order_itemid": "2345"}]';

-------------------------------------------------------------------------------------------------------------------------------
