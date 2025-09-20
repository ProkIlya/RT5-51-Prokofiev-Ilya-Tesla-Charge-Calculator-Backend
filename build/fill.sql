-- Заполнение пользователей
INSERT INTO users (login, password, is_moderator) VALUES
  ('user', 'user', false),
  ('admin', 'admin', true),
  ('user2', 'user2', false);

-- Заполнение условий езды
INSERT INTO driving_scenarios (name, description, status, image_url, type, system_consuption, speed, aero_coeff, rolling_coeff) VALUES
  (
    'Экономичная городская поездка',
    'Поездка по городу без каких-либо систем комфорта',
    'действует',
    'http://localhost:9000/tesla-images/1.jpg',
    'дорога',
    0,
    50,
    0.208,
    0.012
  ),
  (
    'Кондиционер',
    'Работа системы кондиционирования увеличивается потребление энергии',
    'действует',
    'http://localhost:9000/tesla-images/2.jpg',
    'комфорт',
    1.5,
    0,
    0,
    0
  ),
  (
    'Зимний город',
    'Заснежененная трасса и зимняя резина увеличивает трение и ухудшает аэродинамику',
    'действует',
    'http://localhost:9000/tesla-images/3.jpg',
    'дорога',
    3.0,
    50,
    0.228,
    0.018
  ),
  (
    'Загородная трасса',
    'Вождение по шоссе на высоких скоростях',
    'действует',
    'http://localhost:9000/tesla-images/4.jpg',
    'дорога',
    0,
    110,
    0.208,
    0.012
  ),
  (
    'Зимняя загородная трасса',
    'Высокая скорость, заснеженная трасса и потеря аэродинамики',
    'действует',
    'http://localhost:9000/tesla-images/5.jpg',
    'дорога',
    2.0,
    120,
    0.228,
    0.018
  );

-- Заполнение заявок
INSERT INTO trip_applications (status, creator_id, start_charge, created_at) VALUES
  ('черновик', 1, 100, CURRENT_TIMESTAMP);

-- Заполнение связи заявок и условий
INSERT INTO trip_scenarios (trip_application_id, driving_scenario_id, duration) VALUES
  (1, 1, 50),
  (1, 2, 2);