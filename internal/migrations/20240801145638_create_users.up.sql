CREATE TABLE public.users (
                              id uuid NOT NULL,
                              login varchar(100) NOT NULL,
                              "password" varchar(256) NOT NULL,
                              balance numeric(10, 4) DEFAULT 0 NOT NULL,
                              created_date timestamp with time zone NOT NULL,
                              CONSTRAINT users_pk PRIMARY KEY (id),
                              CONSTRAINT users_unique UNIQUE (login)
);
COMMENT ON TABLE public.users IS 'Пользователи';

-- Column comments

COMMENT ON COLUMN public.users.id IS 'Идентификатор пользователя';
COMMENT ON COLUMN public.users.login IS 'Логин';
COMMENT ON COLUMN public.users."password" IS 'Хеш пароля';
COMMENT ON COLUMN public.users.balance IS 'Баланс';
COMMENT ON COLUMN public.users.created_date IS 'Дата создания';
