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
        h3{
            display: flex;
            width: 30%;
            margin-left: auto;
            margin-right: auto;
            margin-top: 20px;
        }
        h1{
            color: #2c3e50;
            display: flex;
            width: 25%;
            margin-left: auto;
            margin-right: auto;
            
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
        .m-desc{
            width: 400px;
            height: 50px;
            max-width: 400px;
            max-height: 50px;
            min-width: 400px;
            min-height: 50px;
            padding: 0.2rem;
            border: 2px solid #b4b4b4;
            border-radius: 10px;
            font-size: 20px;
            box-sizing: border-box;
            transition: all 0.3s ease;
            background: rgba(255, 255, 255, 0.8);
            margin: 10px ;
            margin-top: 20px;
            margin-left: -15px;
        }
        .desc{
            width: 400px;
            height: 200px;
            max-width: 400px;
            max-height: 200px;
            min-width: 400px;
            min-height: 200px;
            padding: 0.2rem;
            border: 2px solid #b4b4b4;
            border-radius: 10px;
            font-size: 20px;
            box-sizing: border-box;
            transition: all 0.3s ease;
            background: rgba(255, 255, 255, 0.8);
            margin: 10px ;
            margin-top: 20px;
            margin-left: -15px;
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
            margin-top: auto;
            
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

        .descriptions{
            margin-top: 50px;
            margin-left: 37.2%;  
            display: flex; 
            flex-direction: column; 
            min-height: 100vh;         
        }
          
        .field__wrapper {
  width: 100%;
  position: relative;
  margin: 15px 0;
  text-align: center;
}
 
.field__file {
  opacity: 0;
  visibility: hidden;
  position: absolute;
}
 
.field__file-wrapper {
  width: 31%;
  display: -webkit-box;
  display: -ms-flexbox;
  display: flex;
  -webkit-box-pack: justify;
      -ms-flex-pack: justify;
          justify-content: space-between;
  -webkit-box-align: center;
      -ms-flex-align: center;
          align-items: center;
  -ms-flex-wrap: wrap;
      flex-wrap: wrap;
}
 
.field__file-fake {
  height: 60px;
  width: calc(100% - 130px);
  display: -webkit-box;
  display: -ms-flexbox;
  display: flex;
  -webkit-box-align: center;
      -ms-flex-align: center;
          align-items: center;
  padding: 0 15px;
  border: 1px solid #c7c7c7;
  border-radius: 3px 0 0 3px;
  border-right: none;
}
 
.field__file-button {
  width: 130px;
  height: 60px;
  background: #3498db;
  color: #fff;
  font-size: 1rem;
  font-weight: 700;
  display: -webkit-box;
  display: -ms-flexbox;
  display: flex;
  -webkit-box-align: center;
      -ms-flex-align: center;
          align-items: center;
  -webkit-box-pack: center;
      -ms-flex-pack: center;
          justify-content: center;
  border-radius: 0 3px 3px 0;
  cursor: pointer;
  transition: background-color 0.3s;
}
.field__file-button:hover {
        background-color: #2980b9;
    }

    .SaveBtn{
        background-color: #3498db;
        color: white;
        border: none;
        padding: 8px 15px;
        border-radius: 5px;
        cursor: pointer;
        margin-top: 25px;
        margin-left: 230px;
        transition: background-color 0.3s;
        width: 150px;
    }
    .SaveBtn:hover {
        background-color: #2980b9;
    }
    .SendBtn{
        background-color: #3498db;
        color: white;
        border: none;
        padding: 8px 10px;
        border-radius: 5px;
        cursor: pointer;
        width: 250px;
        margin-left: 4.5%;
        
        transition: background-color 0.3s;
    }
    .SendBtn:hover {
        background-color: #2980b9;
    }
    </style>
</head>
<body>
    <header class="app-header">
        <div class="header-content">
            <a href="#" class="logo">
                <i class="fas fa-database logo-icon"></i>
                <span>SQL Practice</span>
            </a>
            <nav class="nav-menu">
                <a href="load_offer.php" class="nav_element">Предложить задачу</a>
                <a href="sandbox.php" class="nav_element">Песочница</a>
                <a href="guid.php" class="nav_element">Справочник</a>
                <div class="dropdown">
                    <button class="dropdown-btn">
                        <i class="fas fa-bars menu-icon"></i>
                        Задачи
                    </button>
                    <div class="dropdown-content">
                        <a href="list_tasks.php#easy" class="dropdown-item">Легкие</a>
                        <a href="list_tasks.php#medium" class="dropdown-item">Средние</a>
                        <a href="list_tasks.php#hard" class="dropdown-item">Сложные</a>
                    </div>
                </div>
                <a href="login.php" class="nav_element"><?=isset($_COOKIE["login"])?$_COOKIE["login"]:'Авторизация'?></a>
            </nav>
        </div>
    </header>
    <?php
    if(!isset($_COOKIE['login']))
    { 
    ?>
    <h1>Для того чтобы предложить свою задачу необходимо <a href='login.php'>авторизоваться</a></h1>
    <?php
    }
    else
    {
    ?>
    <div>
        <p>Для отправки своей задачи необходимо заполнить форму расположенную ниже.</p>
        <p>В заголовке задачи необходимо кратко обозначить суть задачи (например: "Нахождение средней скорости движения")</p>
        <p>В описании задачи необходимо подробно изложить суть задачи, описать таблицы и данные, которые в них храняться.</p>
        <p>Далее необходимо прикрепить файлы в формате .csv, в которых храняться данные.</p>
    </div>
    
    <div class="descriptions">
            <h2>Заголовок задачи:</h2>
            <textarea class="m-desc" name="small" id="small"></textarea>
            
        
            <h2>Описание задачи:</h2>
            <textarea class="desc" name="full" id="full"></textarea>
            
        
            <h2>Файл:</h2>
            <div class="field__wrapper">
  
                <input name="file" type="file" name="file" id="field__file-2" class="field field__file" multiple>
                   
                <label class="field__file-wrapper" for="field__file-2">
                  <div class="field__file-fake">Файл не выбран</div>
                  <div class="field__file-button">Добавить файл</div>
                </label>
                   
             </div>
            
            <div class="Send">
                <button class="SendBtn">🚀 Отправить</button>
            </div>
        
    </div>
    
    <?php
    }
    ?>

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
        let fields = document.querySelectorAll('.field__file');
        const filesNames = new Array();
        const filesData = new Array();
        Array.prototype.forEach.call(fields, function (input) {
            let label = input.nextElementSibling,
            labelVal = label.querySelector('.field__file-fake').innerText;
            
            input.addEventListener('change', function (e) {
            
                console.log(this.files);
                for(let i = 0; i < this.files.length; i++)
                {
                    filesNames.push(this.files[i].name);
                    const reader = new FileReader();
        
                    reader.onload = function(e) {
                        filesData.push(e.target.result);
                    };

                    reader.readAsText(this.files[i]);
                }
                
                
                if (filesNames.length > 0)
                label.querySelector('.field__file-fake').innerText = 'Выбрано файлов: ' + filesNames;
                else
                label.querySelector('.field__file-fake').innerText = labelVal;
          });
        });
        var button = document.querySelector('.SendBtn');
        button.addEventListener('click', function(event) {
            let header = document.getElementById('small');
            let descript = document.getElementById('full');
            const formData = {
                'header' : document.getElementById('small').value,
                'description' : document.getElementById('full').value,
                'fileNames' : filesNames,
                'contents' : filesData
            }

            fetch('http://192.168.1.8:8000/addTask', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(formData)
            })
            .then(response => response.json())
            .then(data => {
                if (data.status !== '200 OK') {
                    showError(data.status);
                } else {
                    alert("Данные успешно загружены");
                }
            })
            .catch(error => {
                showError('Request failed: ' + error.message);
            });

        });

    </script>
</body>
</html>
