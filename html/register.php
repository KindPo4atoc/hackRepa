<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Регистрация</title>
    <link href="https://fonts.googleapis.com/css2?family=Poppins:wght@400;500;600;700&display=swap" rel="stylesheet">
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0-beta3/css/all.min.css">
    <link rel="stylesheet" href="../css/register.css">
    <link rel="stylesheet" href="../css/lr.css">
</head>
<body>
    <div class="register-container">
        <h1>Регистрация</h1>
        
        <?php if (!empty($errors['general'])): ?>
            <div class="error-message"><?php echo htmlspecialchars($errors['general']); ?></div>
        <?php endif; ?>
        
        <form id="registerForm" method="POST">
            <div class="form-group">
                <label for="username">Имя пользователя</label>
                <div class="input-with-icon">
                    <i class="fas fa-user"></i>
                    <input placeholder="Введите логин" type="text" id="username" name="username" value="" required>
                </div>
                <?php if (!empty($errors['username'])): ?>
                    <span class="field-error"><?php echo htmlspecialchars($errors['username']); ?></span>
                <?php endif; ?>
            </div>
            
            <div class="form-group">
                <label for="password">Пароль</label>
                <div class="input-with-icon">
                    <i class="fas fa-lock"></i>
                    <input placeholder="Введите пароль" type="password" id="password" name="password" required>
                </div>
                <?php if (!empty($errors['password'])): ?>
                    <span class="field-error"><?php echo htmlspecialchars($errors['password']); ?></span>
                <?php endif; ?>
            </div>
            
            <div class="form-group">
                <label for="confirm_password">Подтвердите пароль</label>
                <div class="input-with-icon">
                    <i class="fas fa-lock"></i>
                    <input placeholder="Введите пароль" type="password" id="confirm_password" name="confirm_password" required>
                </div>
                <?php if (!empty($errors['confirm_password'])): ?>
                    <span class="field-error"><?php echo htmlspecialchars($errors['confirm_password']); ?></span>
                <?php endif; ?>
            </div>
            
            <button type="submit" class="register-button">Зарегистрироваться</button>
        </form>
        
        <div class="links">
            <span>Уже есть аккаунт? <a href="login.php">Войти</a></span>
        </div>
    </div>

    <script>
        document.querySelector('form').addEventListener('submit', function(e) {
            e.preventDefault();
            if(document.getElementById('password').value == document.getElementById('confirm_password').value)
            {
                const formData = {
                    login: document.getElementById('username').value,
                    pass_hash: document.getElementById('password').value
                };
                console.log(formData);
                fetch('http://192.168.1.8:8000/addUser', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(formData)
                })
                .then(response => response.json())
                .then(data => {
                    if (data.status != '200 OK') {
                        alert(data.error);
                    } else {
                        window.location.href = 'login.php';
                    }
                })
                .catch((error) => {
                    console.error('Error:', error);
                    alert('Сревер временно недоступен');
                });
            }
            else
            {
                alert('Пароли не совпадают');
            }
        });
    </script>
</body>
</html>