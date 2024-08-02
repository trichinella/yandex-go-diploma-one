CREATE TABLE public.orders (
                               id uuid NOT NULL,
                               user_id uuid NOT NULL,
                               "number" bigint NOT NULL,
                               created_date timestamp with time zone NOT NULL,
                               status_id uuid NOT NULL,
                               accrual numeric(10, 4) DEFAULT NULL NULL,
                               paid numeric(10, 4) NOT NULL,
                               CONSTRAINT orders_pk PRIMARY KEY (id),
                               CONSTRAINT orders_unique UNIQUE ("number"),
                               CONSTRAINT orders_users_fk FOREIGN KEY (user_id) REFERENCES public.users(id),
                               CONSTRAINT orders_order_statuses_fk FOREIGN KEY (status_id) REFERENCES public.order_statuses(id)
);
COMMENT ON TABLE public.orders IS 'Заказы';

-- Column comments

COMMENT ON COLUMN public.orders.id IS 'Идентификатор заказа';
COMMENT ON COLUMN public.orders.user_id IS 'ID пользователя';
COMMENT ON COLUMN public.orders."number" IS 'Номер заказа';
COMMENT ON COLUMN public.orders.created_date IS 'Дата создания';
COMMENT ON COLUMN public.orders.status_id IS 'Статус заказа';
COMMENT ON COLUMN public.orders.accrual IS 'Баллы лояльности, начисленные за заказ';
COMMENT ON COLUMN public.orders.paid IS 'Часть заказа, оплаченная баллами';
