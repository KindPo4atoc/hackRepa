<?php
   if ($_COOKIE['role'] == '1')
   {
?>
<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>SQL Практика</title>
    <script src="https://smtpjs.com/v3/smtp.js"></script>
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
        h4{
            margin-top: 5px;
        }
        h4, h2, h3 {
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
            margin: 20px;
            
            
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
        .AcceptBtn{
            background-color: #6eff5b93;
            color: rgb(0, 0, 0);
            border: none;
            padding: 8px 15px;
            border-radius: 5px;
            cursor: pointer;
            margin-top: 10px;
            transition: background-color 0.3s;
        }
        .AcceptBtn:hover {
            background-color: #35b929;
        }
        .CancelBtn{
            background-color: #ff5b5b93;
            color: rgb(0, 0, 0);
            border: none;
            padding: 8px 15px;
            border-radius: 5px;
            cursor: pointer;
            margin-top: 10px;
            transition: background-color 0.3s;
        }
        .CancelBtn:hover {
            background-color: #b92929;
        }
        .DownloadBtn{
            background-color: #3498db;
            color: white;
            border: none;
            padding: 8px 15px;
            border-radius: 5px;
            cursor: pointer;
            margin-top: 10px;
            transition: background-color 0.3s;
        }
        .DownloadBtn:hover {
            background-color: #2980b9;
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
        input{
            width: 150px;
            height: 25px;
            padding: 0.2rem;
            border: 2px solid #b4b4b4;
            border-radius: 10px;
            font-size: 15px;
            box-sizing: border-box;
            transition: all 0.3s ease;
            background: rgba(255, 255, 255, 0.8);
            margin-top: 20px ;
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
        .dropdown {
            position: relative;
            display: inline-block;
        }

        .dropdown-btn {
            background-color: #3498db;
            color: white;
            padding: 12px 24px;
            border: none;
            cursor: pointer;
            border-radius: 5px;
            display: flex;
            align-items: center;
            gap: 8px;
            transition: background-color 0.3s;
        }

        .dropdown-btn:hover {
            background-color: #2980b9;
        }

        .dropdown-content {
            display: none;
            position: absolute;
            top: 100%;
            left: 0;
            background-color: #f9f9f9;
            min-width: 100px;
            box-shadow: 0 8px 16px rgba(0,0,0,0.2);
            border-radius: 5px;
            overflow: hidden;
            z-index: 1000;
        }

        .dropdown:hover .dropdown-content {
            display: block;
            animation: slideDown 0.3s ease-out;
        }

        .dropdown-item {
            color: #333;
            background-color:rgb(119, 171, 206);
            padding: 15px 53px;
            text-decoration: none;
            display: block;
            transition: background-color 0.2s;
        }

        .dropdown-item:hover {
            background-color: #f1f1f1;
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
                <a href="ready_tasks.php" class="nav_element">Готовые задачи</a>
                <a href="offer.php" class="nav_element">Предложения</a>
                <div class="dropdown">
                    <button class="dropdown-btn">
                        <i class="fas fa-bars menu-icon"></i>
                        <a href="login.php" class="nav_element"><?=isset($_COOKIE["login"])?$_COOKIE["login"]:'Авторизация'?></a>
                    </button>
                    <div class="dropdown-content">
                        <button class="dropdown-item" id="stoped">Выйти</button>
                    </div>
                </div>
                
            </nav>
        </div>
    </header>

    <div id="tasks">
        
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
        let cancelBtn;
        async function loadTasks() {
            try {
                const response = await fetch(`http://192.168.1.8:8000/getTasksToCheck`);
                if (!response.ok) throw new Error('Ошибка загрузки задач');
                const tasks = await response.json();
                console.log(tasks)
                const container = document.getElementById(`tasks`);
                container.innerHTML = '';
                
                tasks.tasks.forEach(task => {
                    const container_task = document.createElement('div');
                    container_task.className = 'task-container';
                    const card = document.createElement('div');
                    card.className = 'task-card';
                    card.innerHTML = `
                        <h4>От пользователя: ${task.author}</h4>
                        <h3>${task.header}</h3>
                        
                        <button class="show-solution" onclick="toggleSolution('solution${task.id}')">Узнать больше</button>
                        <div id="solution${task.id}" class="solution">
                            <p>Информация:</p>
                            <div class="sql-code">${task.description}</div>
                            <button class="DownloadBtn">Скачать данные</button>
                            <button class="AcceptBtn" id="${task.id}">✔️ Принять</button>
                            <button class="CancelBtn" id="${task.id}">❌ Отклонить</button>
                            <input placeholder="Введите сложность" id="level${task.id}" type="text">
                        </div>
                    `;
                    container.appendChild(card);

                    container.querySelectorAll('.CancelBtn').forEach(btn => {
                    btn.addEventListener('click', async function() {
                        const bntId = this.id;
                        const response = await fetch(`http://192.168.1.8:8000/changeStatus/${bntId}/${0}/s`);
                        if (!response.ok) throw new Error('Ошибка загрузки задач');
                        const tasks = await response.json();
                        console.log(tasks);
                        if(tasks.status == '200 OK')
                            alert('Оклонено успешно');
                        else
                            alert('Ошибка!');
                        location.reload();
                        });
                    });
                    container.querySelectorAll('.AcceptBtn').forEach(btn => {
                    btn.addEventListener('click', async function() {
                        const bntId = this.id;
                        const textBox = document.getElementById('level'+bntId).value;
                        const response = await fetch(`http://192.168.1.8:8000/changeStatus/${bntId}/${1}/${textBox}`);
                        if (!response.ok) throw new Error('Ошибка загрузки задач');
                        const tasks = await response.json();
                        if(tasks.status == '200 OK')
                            alert('Одобрено успешно');
                        else
                            alert('Ошибка!');
                        location.reload();
                    });
                    });
                });
                cancelBtn = document.getElementById('cancel');
            } catch (error) {
                console.error('Ошибка:', error);
                const container = document.getElementById(`${category}-container`);
                container.innerHTML = `<p class="error">Ошибка загрузки задач: ${error.message}</p>`;
            }
        }
        
        function deleteCookie(name) {
            document.cookie = `${name}=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;`;
        }

        let btnStop = document.getElementById('stoped');
        btnStop.addEventListener('click', function(){
            console.log("djnv")
            deleteCookie('login');
            deleteCookie('role');
        });


        // Загрузка задач при старте
        document.addEventListener('DOMContentLoaded', () => {
            loadTasks();
        });

        function toggleSolution(id) {
            const solution = document.getElementById(id);
            const button = solution.previousElementSibling;
            
            if (solution.style.display === 'block') {
                solution.style.display = 'none';
                button.textContent = 'Узнать больше';
            } else {
                solution.style.display = 'block';
                button.textContent = 'Скрыть';
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
<?php
    }
    else
    {
        header('Location: list_tasks.php');
    }
?>

