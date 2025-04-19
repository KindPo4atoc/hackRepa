<?php
session_start();

// Если пользователь уже авторизован, перенаправляем
if (isset($_SESSION['user_id'])) {
    header("Location: login.php");
    exit;
}

$errors = [];
$success = false;
$username = '';
$email = '';

if ($_SERVER['REQUEST_METHOD'] === 'POST') {
    $username = trim($_POST['username'] ?? '');
    $email = trim($_POST['email'] ?? '');
    $password = trim($_POST['password'] ?? '');
    $confirm_password = trim($_POST['confirm_password'] ?? '');

    // Валидация данных
    if (empty($username)) {
        $errors['username'] = 'Имя пользователя обязательно';
    } elseif (strlen($username) < 4) {
        $errors['username'] = 'Не менее 4 символов';
    }

    if (empty($email)) {
        $errors['email'] = 'Email обязателен';
    } elseif (!filter_var($email, FILTER_VALIDATE_EMAIL)) {
        $errors['email'] = 'Некорректный email';
    }

    if (empty($password)) {
        $errors['password'] = 'Пароль обязателен';
    } elseif (strlen($password) < 6) {
        $errors['password'] = 'Не менее 6 символов';
    }

    if ($password !== $confirm_password) {
        $errors['confirm_password'] = 'Пароли не совпадают';
    }

    // Если нет ошибок - отправка на API
    if (empty($errors)) {
        $api_url = 'http://localhost:5000/api/register';
        $data = [
            'username' => $username,
            'email' => $email,
            'password' => $password
        ];

        
        $url = 'http://example.com/api';
$options = [
    'http' => [
        'method' => 'GET', 
        'header' => "User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64)\r\n"
    ]
];

$context = stream_context_create($options);
$response = file_get_contents($url, false, $context);

if ($response === false) {
    die("Ошибка запроса: " . print_r(error_get_last(), true));
}

echo $response;


        $response = curl_exec($ch);
        $http_code = curl_getinfo($ch, CURLINFO_HTTP_CODE);
        $curl_error = curl_error($ch);
        curl_close($ch);

        // Обработка ответа
        if ($response === false) {
            $errors['general'] = 'Ошибка подключения к API: ' . $curl_error;
        } else {
            $result = json_decode($response, true);
            
            if (json_last_error() !== JSON_ERROR_NONE) {
                $errors['general'] = 'Неверный формат ответа от сервера';
            } elseif ($http_code === 201 && isset($result['message'])) {
                $_SESSION['registration_success'] = $result['message'];
                header("Location: login.php");
                exit;
            } else {
                $errors['general'] = $result['error'] ?? 'Ошибка регистрации';
            }
        }
    }
}
?>
<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Регистрация</title>
    <link href="https://fonts.googleapis.com/css2?family=Poppins:wght@400;500;600;700&display=swap" rel="stylesheet">
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0-beta3/css/all.min.css">
    <link rel="stylesheet" href="register.css">
    <link rel="stylesheet" href="main.css">
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
                    <input type="text" id="username" name="username" value="<?php echo htmlspecialchars($username); ?>" required>
                </div>
                <?php if (!empty($errors['username'])): ?>
                    <span class="field-error"><?php echo htmlspecialchars($errors['username']); ?></span>
                <?php endif; ?>
            </div>
            
            <div class="form-group">
                <label for="email">Email</label>
                <div class="input-with-icon">
                    <i class="fas fa-envelope"></i>
                    <input type="email" id="email" name="email" value="<?php echo htmlspecialchars($email); ?>" required>
                </div>
                <?php if (!empty($errors['email'])): ?>
                    <span class="field-error"><?php echo htmlspecialchars($errors['email']); ?></span>
                <?php endif; ?>
            </div>
            
            <div class="form-group">
                <label for="password">Пароль</label>
                <div class="input-with-icon">
                    <i class="fas fa-lock"></i>
                    <input type="password" id="password" name="password" required>
                </div>
                <?php if (!empty($errors['password'])): ?>
                    <span class="field-error"><?php echo htmlspecialchars($errors['password']); ?></span>
                <?php endif; ?>
            </div>
            
            <div class="form-group">
                <label for="confirm_password">Подтвердите пароль</label>
                <div class="input-with-icon">
                    <i class="fas fa-lock"></i>
                    <input type="password" id="confirm_password" name="confirm_password" required>
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
        const express = require('express');
const cors = require('cors'); 

const app = express();


app.use(cors());


app.use(cors({
  origin: 'http://localhost:3000',
  methods: ['GET', 'POST', 'PUT', 'DELETE'],
  allowedHeaders: ['Content-Type', 'Authorization']
}));


app.post('/api/register', (req, res) => {
  res.json({ success: true });
});

app.listen(5000, () => console.log('Server running on port 5000'));

        document.getElementById('registerForm').addEventListener('submit', async function(e) {
            e.preventDefault();
            
            const formData = {
                username: document.getElementById('username').value,
                email: document.getElementById('email').value,
                password: document.getElementById('password').value,
                confirm_password: document.getElementById('confirm_password').value
            };

            try {
                const response = await fetch('http://localhost:5000/api/register', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(formData)
                });

                
                const contentType = response.headers.get('content-type');
                if (!contentType || !contentType.includes('application/json')) {
                    const text = await response.text();
                    throw new Error(`Ожидался JSON, получили: ${text.substring(0, 100)}...`);
                }

                const data = await response.json();

                if (!response.ok) {
                    throw new Error(data.error || 'Ошибка сервера');
                }

                
                window.location.href = 'login.php?registration=success';
            } catch (error) {
                console.error('Registration Error:', error);
                alert(error.message || 'Ошибка при регистрации');
            }
        });
    </script>
</body>
</html>