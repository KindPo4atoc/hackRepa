<?php
session_start();

// Редирект если уже авторизован
if (isset($_SESSION['user_id'])) {
    header("Location: list_tasks.php");
    exit;
}
?>
<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Вход в систему</title>
    <link href="https://fonts.googleapis.com/css2?family=Poppins:wght@400;500;600;700&display=swap" rel="stylesheet">
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0-beta3/css/all.min.css">
    <link rel="stylesheet" href="login.css">
    <link rel="stylesheet" href="main.css">
</head>
<body>
    <div class="login-container">
        <h1>Вход в систему</h1>
        
        <?php if (!empty($error)): ?>
            <div class="error-message"><?php echo htmlspecialchars($error); ?></div>
        <?php endif; ?>
        
        <form action="login.php" method="POST">
            <div class="form-group">
                <label for="username">Логин</label>  <!-- Изменили с Email на Логин -->
                <div class="input-with-icon">
                    <i class="fas fa-user"></i>  <!-- Изменили иконку -->
                    <input type="text" id="username" name="username" value="<?php echo htmlspecialchars($username); ?>" required>
                </div>
            </div>
            
            <div class="form-group">
                <label for="password">Пароль</label>
                <div class="input-with-icon">
                    <i class="fas fa-lock"></i>
                    <input type="password" id="password" name="password" required>
                </div>
            </div>
            
            <button type="submit" class="login-button">Войти</button>
        </form>
        
        <div class="links">
            <a href="register.php">Регистрация</a>
            <a href="forgot_password.php">Забыли пароль?</a>
        </div>
    </div>

    <script>
        // AJAX версия отправки формы
        document.querySelector('form').addEventListener('submit', function(e) {
            e.preventDefault();
            
            const formData = {
                username: document.getElementById('username').value,
                password: document.getElementById('password').value
            };
            
            fetch('http://localhost:5000/api/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(formData)
            })
            .then(response => response.json())
            .then(data => {
                if (data.error) {
                    alert(data.error);
                } else {
                    // Перенаправляем на защищенную страницу
                    window.location.href = 'dashboard.php';
                }
            })
            .catch((error) => {
                console.error('Error:', error);
                alert('Произошла ошибка при входе');
            });
        });
    </script>
</body>
</html>