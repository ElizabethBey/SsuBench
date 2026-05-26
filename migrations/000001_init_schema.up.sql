-- Типы ролей
CREATE TYPE user_role AS ENUM ('customer', 'executor', 'admin');

-- Статусы задач
CREATE TYPE task_status AS ENUM ('open', 'in_progress', 'completed', 'cancelled');

-- Статусы откликов
CREATE TYPE bid_status AS ENUM ('pending', 'accepted', 'rejected');

-- Таблица пользователей
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role user_role NOT NULL,
    balance DECIMAL(15, 2) DEFAULT 0.00 NOT NULL,
    is_blocked BOOLEAN DEFAULT FALSE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Таблица задач
CREATE TABLE tasks (
   id SERIAL PRIMARY KEY,
   customer_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
   title VARCHAR(255) NOT NULL,
   description TEXT NOT NULL,
   budget DECIMAL(15, 2) NOT NULL,
   status task_status DEFAULT 'open' NOT NULL,
   created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Таблица откликов
CREATE TABLE bids (
  id SERIAL PRIMARY KEY,
  task_id INTEGER REFERENCES tasks(id) ON DELETE CASCADE,
  executor_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
  amount DECIMAL(15, 2) NOT NULL,
  status bid_status DEFAULT 'pending' NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
-- Ограничение: исполнитель может оставить только один отклик на одну задачу
  UNIQUE(task_id, executor_id)
);

-- Таблица платежей
CREATE TABLE payments (
  id SERIAL PRIMARY KEY,
  task_id INTEGER REFERENCES tasks(id),
  sender_id INTEGER REFERENCES users(id),
  receiver_id INTEGER REFERENCES users(id),
  amount DECIMAL(15, 2) NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_bids_task ON bids(task_id);
