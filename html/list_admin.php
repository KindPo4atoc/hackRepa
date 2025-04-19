<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>SQL Практика</title>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0/css/all.min.css">
    <style>
        * {
            box-sizing: border-box;
            margin: 0;
            padding: 0;
        }

        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            line-height: 1.6;
            background-color: #f5f7fa;
            color: #2c3e50;
        }
        
        h1, h2, h3 {
            color: #2c3e50;
        }
        
        /* Стили для навигации */
        .nav {
            display: flex;
            background-color: #3498db;
            padding: 10px;
            border-radius: 5px;
            margin-bottom: 20px;
        }
        
        .nav a {
            color: white;
            text-decoration: none;
            padding: 10px 15px;
            margin-right: 10px;
            border-radius: 3px;
            transition: background-color 0.3s;
        }
        
        .nav a:hover {
            background-color: #2980b9;
        }
        
        /* Стили для задач */
        .task-container {
            display: flex;
            flex-wrap: wrap;
            gap: 20px;
        }
        
        .task-card {
            background-color: white;
            border-radius: 8px;
            box-shadow: 0 2px 5px rgba(0,0,0,0.1);
            padding: 20px;
            flex: 1 1 300px;
            transition: transform 0.3s;
        }
        
        .task-card:hover {
            transform: translateY(-5px);
            box-shadow: 0 5px 15px rgba(0,0,0,0.1);
        }
        
        .task-card h3 {
            margin-top: 0;
            border-bottom: 2px solid #3498db;
            padding-bottom: 10px;
        }
        
        .difficulty {
            display: inline-block;
            padding: 3px 8px;
            border-radius: 3px;
            font-size: 0.8em;
            margin-bottom: 10px;
        }
        
        .easy {
            background-color: #d4edda;
            color: #155724;
        }
        
        .medium {
            background-color: #fff3cd;
            color: #856404;
        }
        
        .hard {
            background-color: #f8d7da;
            color: #721c24;
        }
        
        /* Стили для SQL кода */
        .sql-code {
            background-color: #f0f0f0;
            border-left: 4px solid #3498db;
            padding: 15px;
            margin: 15px 0;
            font-family: 'Courier New', Courier, monospace;
            white-space: pre-wrap;
            border-radius: 0 5px 5px 0;
            overflow-x: auto;
        }
        
        /* Стили для решения */
        .solution {
            display: none;
            margin-top: 15px;
            padding: 15px;
            background-color: #e8f4fc;
            border-radius: 5px;
        }
        
        .show-solution {
            background-color: #3498db;
            color: white;
            border: none;
            padding: 8px 15px;
            border-radius: 5px;
            cursor: pointer;
            margin-top: 10px;
            transition: background-color 0.3s;
        }
        
        .show-solution:hover {
            background-color: #2980b9;
        }
        
        /* Стили для фильтров */
        .filters {
            background-color: white;
            padding: 15px;
            border-radius: 8px;
            margin-bottom: 20px;
            box-shadow: 0 2px 5px rgba(0,0,0,0.1);
        }
        
        .filter-group {
            margin-bottom: 10px;
        }
        
        .filter-group label {
            margin-right: 10px;
        }
        
        /* Адаптивность */
        @media (max-width: 768px) {
            .task-card {
                flex: 1 1 100%;
            }
            
            .nav {
                flex-direction: column;
            }
            
            .nav a {
                margin-bottom: 5px;
                margin-right: 0;
            }
        }

        .app-header {
            background-color: #2c3e50;
            padding: 1rem 2rem;
            position: sticky;
            top: 0;
            z-index: 1000;
            box-shadow: 0 2px 5px rgba(0,0,0,0.1);
        }

        .header-content {
            max-width: 1200px;
            margin: 0 auto;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }

        .logo {
            display: flex;
            align-items: center;
            color: white;
            text-decoration: none;
            gap: 10px;
        }

        .logo-icon {
            font-size: 1.5rem;
            color: #3498db;
        }

        .nav-menu {
            display: flex;
            gap: 20px;
        }

        .nav-link {
            color: white;
            text-decoration: none;
            padding: 0.5rem 1rem;
            border-radius: 4px;
            transition: background-color 0.3s;
        }

        .nav-link:hover {
            background-color: #34495e;
        }
        .nav_element{
            color: white;
            text-decoration: none;
            padding: 0.5rem 1rem;
            border-radius: 4px;
            transition: background-color 0.3s;
        }

        .app-footer {
            background-color: #2c3e50;
            color: white;
            padding: 2rem;
            margin-top: 40px;
        }

        .footer-content {
            max-width: 1200px;
            margin: 0 auto;
            display: flex;
            justify-content: space-between;
            flex-wrap: wrap;
            gap: 20px;
        }

        .footer-section {
            flex: 1;
            min-width: 250px;
        }
        @keyframes slideDown {
            from {
                opacity: 0;
                transform: translateY(-10px);
            }
            to {
                opacity: 1;
                transform: translateY(0);
            }
        }

    </style>
</head>
<body>
    <header class="app-header">
        <div class="header-content">
            <a href="list_admin.php" class="logo">
                <i class="fas fa-database logo-icon"></i>
                <span>SQL Practice</span>
            </a>
            <nav class="nav-menu">
            <a href="ready_tasks.html" class="nav_element">Готовые задачи</a>
            <a href="offer.html" class="nav_element">Предложения</a>
            </nav>
        </div>
    </header>

    <h2 id="easy">Лёгкие задачи</h2>
    <div class="task-container">
        <div class="task-card">
            <span class="difficulty easy">Легко</span>
            <h3>Выбор всех данных из таблицы</h3>
            <p>Напишите SQL запрос для выбора всех записей из таблицы "employees".</p>
            
            <button class="show-solution" onclick="toggleSolution('solution1')">Перейти к решению</button>
            <div id="solution1" class="solution">
                <p>Решение:</p>
                <div class="sql-code">SELECT * FROM employees;</div>
            </div>
        </div>
        
        <div class="task-card">
            <span class="difficulty easy">Легко</span>
            <h3>Выбор конкретных столбцов</h3>
            <p>Напишите SQL запрос для выбора только имен и зарплат сотрудников из таблицы "employees".</p>
            
            <button class="show-solution" onclick="toggleSolution('solution2')">Показать решение</button>
            <div id="solution2" class="solution">
                <p>Решение:</p>
                <div class="sql-code">SELECT first_name, last_name, salary 
FROM employees;</div>
            </div>
        </div>
    </div>
    
    <h2 id="mudium">Средние задачи</h2>
    <div class="task-container">
        <div class="task-card">
            <span class="difficulty medium">Средне</span>
            <h3>INNER JOIN двух таблиц</h3>
            <p>Напишите SQL запрос, который соединяет таблицы "employees" и "departments" по полю department_id и выводит имя сотрудника и название отдела.</p>
            
            <button class="show-solution" onclick="toggleSolution('solution3')">Показать решение</button>
            <div id="solution3" class="solution">
                <p>Решение:</p>
                <div class="sql-code">SELECT e.first_name, e.last_name, d.department_name
FROM employees e
INNER JOIN departments d ON e.department_id = d.department_id;</div>
            </div>
        </div>
        
        <div class="task-card">
            <span class="difficulty medium">Средне</span>
            <h3>Подзапрос в WHERE</h3>
            <p>Напишите SQL запрос для вывода сотрудников, чья зарплата выше средней по компании.</p>
            
            <button class="show-solution" onclick="toggleSolution('solution5')">Показать решение</button>
            <div id="solution5" class="solution">
                <p>Решение:</p>
                <div class="sql-code">SELECT first_name, last_name, salary
FROM employees
WHERE salary > (SELECT AVG(salary) FROM employees);</div>
            </div>
        </div>
    </div>
    
    <h2 id="subqueries">Сложные задачи</h2> 
    <div class="task-container">
    <div class="task-card">
            <span class="difficulty hard">Сложно</span>
            <h3>LEFT JOIN с агрегацией</h3>
            <p>Напишите SQL запрос, который выводит название отдела и количество сотрудников в каждом отделе, включая отделы без сотрудников.</p>
            
            <button class="show-solution" onclick="toggleSolution('solution4')">Показать решение</button>
            <div id="solution4" class="solution">
                <p>Решение:</p>
                <div class="sql-code">SELECT d.department_name, COUNT(e.employee_id) AS employee_count
FROM departments d
LEFT JOIN employees e ON d.department_id = e.department_id
GROUP BY d.department_name
ORDER BY employee_count DESC;</div>
            </div>
        </div>
    </div>

    <footer class="app-footer">
        <div class="footer-content">
            <div class="footer-section">
                <h4>О проекте</h4>
                <p>Практическая платформа для изучения SQL через решение реальных задач</p>
            </div>
            <div class="footer-section">
                <h4>Контакты</h4>
                <ul>
                    <li>Email: support@sqlpractice.com</li>
                    <li>Телефон: +7 (495) 123-45-67</li>
                </ul>
            </div>
            <div class="footer-section">
                <h4>Социальные сети</h4>
                <div class="social-links">
                    <a href="#" class="nav-link"><i class="fab fa-vk"></i> VK</a>
                    <a href="#" class="nav-link"><i class="fab fa-telegram"></i> Telegram</a>
                </div>
            </div>
        </div>
        <div class="footer-copyright" style="text-align: center; margin-top: 20px;">
            © 2023 SQL Practice. Все права защищены.
        </div>
    </footer>

    <script>
        function toggleSolution(id) {
            const solution = document.getElementById(id);
            const button = solution.previousElementSibling;
            
            if (solution.style.display === 'block') {
                solution.style.display = 'none';
                button.textContent = 'Показать решение';
            } else {
                solution.style.display = 'block';
                button.textContent = 'Скрыть решение';
            }
        }
        document.addEventListener('click', function(event) {
            const dropdowns = document.querySelectorAll('.dropdown');
            dropdowns.forEach(dropdown => {
                if (!dropdown.contains(event.target)) {
                    dropdown.querySelector('.dropdown-content').style.display = 'none';
                }
            });
        });
    </script>
</body>
</html>