<html>
<head>
    <title>Register Upload Request</title>
</head>
<body>
<form enctype="multipart/form-data" action="http://127.0.0.1:8080/register" method="post">
    Your Name:
    <input type="text" name="firstname"><br>
    Case Number:
    <input type="text" name="caseNumber"><br>
    Reason for request:
    <input type="text" name="reason"><br>
    Expires in:
    <input type="number" name="timeValue">
    <input type="radio" name="timeMeasure" value="minute" checked> minutes
    <input type="radio" name="timeMeasure" value="hour"> hours
    <input type="radio" name="timeMeasure" value="day"> days<br><br>
    <input type="submit" value="register" />
</form>
</body>
</html>