CREATE TABLE public.order_statuses (
                                       id uuid NOT NULL,
                                       title varchar(250) NOT NULL,
                                       code varchar(20) NOT NULL,
                                       CONSTRAINT order_statuses_pk PRIMARY KEY (id),
                                       CONSTRAINT order_statuses_unique UNIQUE (code)
);
COMMENT ON TABLE public.order_statuses IS 'Статусы заказа';

-- Column comments

COMMENT ON COLUMN public.order_statuses.id IS 'Идентификатор';
COMMENT ON COLUMN public.order_statuses.title IS 'Заголовок';
COMMENT ON COLUMN public.order_statuses.code IS 'Код';
