const courseButtons = document.querySelectorAll('.course-button');
const courseDescription = document.getElementById('course-description');

courseButtons.forEach((button) => {
    button.addEventListener('click', function () {
        const courseName = this.dataset.course;
        courseDescription.textContent = description;
    });
});