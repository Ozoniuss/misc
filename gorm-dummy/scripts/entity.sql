DROP TABLE IF EXISTS entities;

CREATE TABLE IF NOT EXISTS entities(
    uuid UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    other_name TEXT NOT NULL,
    age INT NOT NULL,
    salary INT NOT NULL
);

INSERT INTO entities (uuid, name, other_name, age, salary) VALUES
('bdbcaa1e-4533-4792-b4b3-4b36a51126c3', 'A', 'B', 1, 4000),
('0edecd47-bcf3-460d-be10-0e7def209c81', 'B', 'E', 2, 5000),
('9bdb2e70-15fb-4683-b59e-3449ac654d63', 'E', 'C', 9, 3000),
('5f875d48-c925-40e7-96ae-2d7ee48e685f', 'F', 'F', 6, 3000),
('cd6e3f87-5c34-44c5-9309-557adf8fba5f', 'B', 'A', 4, 4000),
('5d448886-8e52-4644-9ea2-cf2f7b99f4d7', 'E', 'E', 5, 3000),
('b3097f8d-57dd-4029-a9d7-4fbc7c011bdc', 'A', 'D', 8, 6000),
('acf0e3a7-48c5-4757-bfd7-313b6b249c61', 'B', 'F', 1, 4000),
('7850c63a-95e1-4c93-807c-4c9ee45e974a', 'A', 'F', 8, 5000),
('e4735ab8-0f4e-4d11-ab68-28346d19dd62', 'C', 'A', 7, 3000),
('8463e93f-1431-4996-99ff-db248be041f6', 'C', 'C', 2, 4000),
('228be482-5c6b-48d5-8ea8-1c0a80757079', 'D', 'E', 6, 4000),
('5ba21827-7bd1-46c1-b2f0-b8b1022766b7', 'E', 'B', 9, 3000),
('0b402f7f-bc58-4098-9b56-996f7bca92e7', 'F', 'F', 3, 5000),
('a79000ae-68a7-415e-911c-05d8c00ad92f', 'A', 'A', 5, 3000),
('3be67cd8-476c-4e4d-8f64-141ff3ec1ea1', 'F', 'C', 7, 4000),
('a1c85da1-5e74-46f9-becb-9185e86dceee', 'D', 'D', 2, 3000),
('188868ee-674b-4d52-bdcb-327ad7475dbb', 'C', 'D', 5, 6000),
('c0932bed-e4bb-472d-9567-b648b101abb8', 'D', 'E', 9, 3000),
('b2980af4-0ae3-4346-b0bc-b56ad7370fcd', 'B', 'A', 4, 3000);