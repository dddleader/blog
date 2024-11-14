let currentPage = 1;
const pageSize = 10;

document.addEventListener('DOMContentLoaded', function() {
    // 导航栏活跃状态
    const navLinks = document.querySelectorAll('.nav-links a');
    const currentPath = window.location.pathname;

    navLinks.forEach(link => {
        if (link.getAttribute('href') === currentPath) {
            link.classList.add('active');
        }
    });

    // 加载文章列表
    loadArticles();
});

async function loadArticles() {
    try {
        const response = await fetch(`/api/articles?page=${currentPage}&pageSize=${pageSize}`);
        const data = await response.json();
        
        renderArticles(data.Articles);
        renderPagination(data.Total, data.Page, data.PageSize);
    } catch (error) {
        console.error('Error loading articles:', error);
    }
}

function renderArticles(articles) {
    const grid = document.querySelector('.articles-grid');
    grid.innerHTML = articles.map(article => `
        <article class="article-card" onclick="showArticle('${article.Path}')">
            <img src="${article.Cover || '/static/images/default-cover.jpg'}" alt="${article.Title}" class="article-image">
            <div class="article-content">
                <h3 class="article-title">${article.Title}</h3>
                <p class="article-date">${article.Date}</p>
                <p class="article-excerpt">${article.Summary}</p>
                <div class="article-tags">
                    ${article.Tags.map(tag => `<span class="tag">${tag}</span>`).join('')}
                </div>
            </div>
        </article>
    `).join('');
}

function renderPagination(total, currentPage, pageSize) {
    const totalPages = Math.ceil(total / pageSize);
    const pagination = document.createElement('div');
    pagination.className = 'pagination';
    
    let paginationHTML = '';
    
    if (currentPage > 1) {
        paginationHTML += `<button onclick="changePage(${currentPage - 1})">上一页</button>`;
    }
    
    for (let i = 1; i <= totalPages; i++) {
        paginationHTML += `
            <button class="${i === currentPage ? 'active' : ''}" 
                    onclick="changePage(${i})">${i}</button>
        `;
    }
    
    if (currentPage < totalPages) {
        paginationHTML += `<button onclick="changePage(${currentPage + 1})">下一页</button>`;
    }
    
    pagination.innerHTML = paginationHTML;
    document.querySelector('.articles-grid').after(pagination);
}

async function showArticle(path) {
    try {
        const response = await fetch(`/api/article/${path}`);
        const article = await response.json();
        
        // 创建文章详情模态框
        const modal = document.createElement('div');
        modal.className = 'article-modal';
        modal.innerHTML = `
            <div class="article-modal-content">
                <button class="close-button" onclick="closeArticle()">×</button>
                <h1>${article.Title}</h1>
                <div class="article-meta">
                    <span>${article.Date}</span>
                    <div class="tags">
                        ${article.Tags.map(tag => `<span class="tag">${tag}</span>`).join('')}
                    </div>
                </div>
                <div class="article-body">
                    ${marked(article.Content)}
                </div>
            </div>
        `;
        
        document.body.appendChild(modal);
        document.body.style.overflow = 'hidden';
    } catch (error) {
        console.error('Error loading article:', error);
    }
}

function closeArticle() {
    const modal = document.querySelector('.article-modal');
    if (modal) {
        modal.remove();
        document.body.style.overflow = 'auto';
    }
}

function changePage(page) {
    currentPage = page;
    loadArticles();
    window.scrollTo(0, 0);
}

// 平滑滚动
document.querySelectorAll('a[href^="#"]').forEach(anchor => {
    anchor.addEventListener('click', function (e) {
        e.preventDefault();
        document.querySelector(this.getAttribute('href')).scrollIntoView({
            behavior: 'smooth'
        });
    });
}); 