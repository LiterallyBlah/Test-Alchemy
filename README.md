# 🧪 Test Alchemy

**Transforming inputs into gold-standard test cases.**

Test Alchemy is a lightweight microservice that streamlines the creation, management, and export of user test cases. Whether you prefer manual control or AI assistance, Test Alchemy provides the flexibility and simplicity you need to manage test cases effectively.

---

## ✨ Features

- **Project-Based Organisation**: Create and manage test cases within distinct projects.
- **Two Test Case Creation Methods**:
  - **Manual Input**: Full control over test case details.
  - **AI Assistance**: Generate test cases from images and brief descriptions, or use AI to enhance specific fields.
- **Collaboration**: Invite and manage collaborators within projects.
- **Search & Filter**: Easily find test cases using tags, titles, or user levels.
- **Export**: Export test cases to CSV or Excel formats for easy sharing and reporting.

---

## 🚀 Getting Started

### **Prerequisites**

- **Golang** (Backend)
- **Redis** (Session Management)
- **MySQL** (Database)
- **Temple** (Template Engine)
- **Tailwind CSS** (Styling)

### **Installation**

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/yourusername/test-alchemy.git
   cd test-alchemy
   ```

2. **Set Up the Database**:
   - Configure your MySQL instance and create the necessary tables.

3. **Configure Environment Variables**:
   Create a `.env` file for database and Redis configurations.

4. **Run the Application**:
   ```bash
   go run main.go
   ```

5. **Access the Platform**:
   - Open your browser and navigate to `http://localhost:8080`.

---

## 🧪 **Creating a Test Case**

### **Manual with AI Assistance**
1. Fill in each field manually.
2. Optionally use the **“Generate”** button next to each field to enhance it with AI.

### **Full AI Generation**
1. Upload an image and provide a brief description.
2. Click **“Generate Full Test Case”** to let AI create the test case for you.

---

## 👥 **Collaboration**

- Invite users to your projects via email.
- Manage collaborators through the **Users Modal**.
- Project owners have full control over adding/removing collaborators.

---

## 📦 **Exporting Test Cases**

- Export your test cases as **CSV** or **Excel** files directly from the dashboard.

---

## 🛠️ **Tech Stack**

- **Backend**: Golang
- **Session Management**: Redis
- **Database**: MySQL
- **Frontend**: Temple (Template Engine) + Tailwind CSS

---

## 📄 **Licence**

This project is licensed under the [MIT Licence](LICENSE).

---

## ⭐ **Contributing**

Contributions are welcome! Feel free to open an issue or submit a pull request.

---

## 🌐 **Links**

- **GitHub**: [https://github.com/yourusername/test-alchemy](https://github.com/yourusername/test-alchemy)
