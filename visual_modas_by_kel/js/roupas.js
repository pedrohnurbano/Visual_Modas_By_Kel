// Variável de controle de login
let isUserLoggedIn = localStorage.getItem('userLoggedIn') === 'true';

// DADOS DOS PRODUTOS
const allProducts = [
    {
        id: 1,
        name: "Vestido Midi Floral",
        price: 899.90,
        installments: 7,
        image1: "design/ex-roupa1.png",
        image2: "design/ex-roupa1b.png",
        sizes: ["PP", "P", "M", "G", "GG"],
        category: "vestidos"
    },
    {
        id: 2,
        name: "Blusa Cropped Básica",
        price: 449.90,
        installments: 4,
        image1: "design/ex-roupa2.png",
        image2: "design/ex-roupa2b.png",
        sizes: ["PP", "P", "M", "G", "GG"],
        category: "blusas"
    },
    {
        id: 3,
        name: "Calça Wide Leg Alfaiataria",
        price: 679.90,
        installments: 5,
        image1: "design/ex-roupa3.png",
        image2: "design/ex-roupa3b.png",
        sizes: ["PP", "P", "M", "G", "GG"],
        category: "calcas"
    },
    {
        id: 4,
        name: "Saia Midi Plissada",
        price: 559.90,
        installments: 5,
        image1: "design/ex-roupa4.png",
        image2: "design/ex-roupa4b.png",
        sizes: ["PP", "P", "M", "G", "GG"],
        category: "saias"
    },
    {
        id: 5,
        name: "Shorts Cintura Alta",
        price: 399.90,
        installments: 4,
        image1: "design/ex-roupa5.png",
        image2: "design/ex-roupa5b.png",
        sizes: ["PP", "P", "M", "G", "GG"],
        category: "shorts"
    },
    {
        id: 6,
        name: "Jaqueta Jeans Oversized",
        price: 799.90,
        installments: 6,
        image1: "design/ex-roupa6.png",
        image2: "design/ex-roupa6b.png",
        sizes: ["PP", "P", "M", "G", "GG"],
        category: "jaquetas"
    },
    {
        id: 7,
        name: "Macacão Longo Elegante",
        price: 989.90,
        installments: 8,
        image1: "design/ex-roupa7.png",
        image2: "design/ex-roupa7b.png",
        sizes: ["PP", "P", "M", "G", "GG"],
        category: "macacoes"
    },
    {
        id: 8,
        name: "Blazer Alfaiataria Premium",
        price: 849.90,
        installments: 7,
        image1: "design/ex-roupa8.png",
        image2: "design/ex-roupa8b.png",
        sizes: ["PP", "P", "M", "G", "GG"],
        category: "blazers"
    },
    {
        id: 9,
        name: "Body Canelado Básico",
        price: 329.90,
        installments: 3,
        image1: "design/ex-roupa9.png",
        image2: "design/ex-roupa9b.png",
        sizes: ["PP", "P", "M", "G", "GG"],
        category: "body"
    },
    {
        id: 10,
        name: "Regata Ribana Alongada",
        price: 279.90,
        installments: 3,
        image1: "design/ex-roupa10.png",
        image2: "design/ex-roupa10b.png",
        sizes: ["PP", "P", "M", "G", "GG"],
        category: "regatas"
    },
    {
        id: 11,
        name: "Vestido Longo Festa",
        price: 1199.90,
        installments: 10,
        image1: "design/ex-roupa11.png",
        image2: "design/ex-roupa11b.png",
        sizes: ["PP", "P", "M", "G", "GG"],
        category: "vestidos"
    },
    {
        id: 12,
        name: "Calça Cargo Moderna",
        price: 629.90,
        installments: 5,
        image1: "design/ex-roupa12.png",
        image2: "design/ex-roupa12b.png",
        sizes: ["PP", "P", "M", "G", "GG"],
        category: "calcas"
    }
];

let favorites = JSON.parse(localStorage.getItem('favorites')) || [];
let cart = JSON.parse(localStorage.getItem('cart')) || [];
let selectedSizes = {};
let filteredProducts = [...allProducts];

// FUNÇÕES DE FAVORITOS
function toggleFavorite(productId) {
    if (!isUserLoggedIn) {
        showFavoriteLoginModal();
        return;
    }

    const index = favorites.indexOf(productId);
    if (index > -1) {
        favorites.splice(index, 1);
        showNotification('Produto removido dos favoritos');
    } else {
        favorites.push(productId);
        showNotification('Produto adicionado aos favoritos');
    }
    localStorage.setItem('favorites', JSON.stringify(favorites));
    updateFavoritesCount();
    renderProducts();
}

function showFavoriteLoginModal() {
    const modal = document.getElementById('favoriteLoginModal');
    modal.classList.add('active');
    document.body.style.overflow = 'hidden';
}

function closeFavoriteModal() {
    const modal = document.getElementById('favoriteLoginModal');
    modal.classList.remove('active');
    document.body.style.overflow = '';
}

function updateFavoritesCount() {
    const count = document.getElementById('favoritosCount');
    if (favorites.length > 0) {
        count.textContent = favorites.length;
        count.style.display = 'inline';
    } else {
        count.style.display = 'none';
    }
}

// FUNÇÕES DE SACOLA
function addToCart(productId, size) {
    const product = allProducts.find(p => p.id === productId);
    if (!product) return;

    const existingItem = cart.find(item => item.id === productId && item.size === size);
    if (existingItem) {
        existingItem.quantity++;
    } else {
        cart.push({
            id: productId,
            name: product.name,
            price: product.price,
            size: size,
            quantity: 1,
            image: product.image1
        });
    }
    localStorage.setItem('cart', JSON.stringify(cart));
    updateCartCount();
    showNotification('Produto adicionado à sacola!');
}

function updateCartCount() {
    const count = document.getElementById('sacolaCount');
    const total = cart.reduce((sum, item) => sum + item.quantity, 0);
    if (total > 0) {
        count.textContent = total;
        count.style.display = 'inline';
    } else {
        count.style.display = 'none';
    }
}

// FUNÇÕES DE PRODUTOS
function createProductCard(product) {
    const isFavorite = favorites.includes(product.id);
    const installmentValue = (product.price / product.installments).toFixed(2);

    return `
        <div class="produto-card" data-product-id="${product.id}">
            <button class="produto-favorito ${isFavorite ? 'favorito-active' : ''}" title="${isFavorite ? 'Remover dos favoritos' : 'Adicionar aos favoritos'}" onclick="toggleFavorite(${product.id})">
                <svg width="22" height="22" viewBox="0 0 24 24" fill="none">
                    <path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41 0.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"
                        stroke="#370400" stroke-width="1.5" fill="${isFavorite ? '#370400' : 'none'}"/>
                </svg>
            </button>
            <div class="produto-img-container">
                <img src="${product.image1}" alt="${product.name}" class="produto-img img-primaria">
                <img src="${product.image2}" alt="${product.name} - Verso" class="produto-img img-secundaria">
                <div class="produto-tamanhos">
                    ${product.sizes.map(size => `<span class="size-option" data-size="${size}" onclick="selectSize(${product.id}, '${size}')">${size}</span>`).join('')}
                    <button class="produto-sacola" title="Adicionar à sacola" onclick="addSelectedToCart(${product.id})">
                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="white" stroke-width="2">
                            <path d="M6 2L3 6v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V6l-3-4z"/>
                            <line x1="3" y1="6" x2="21" y2="6"/>
                        </svg>
                    </button>
                </div>
            </div>
            <div class="produto-info">
                <span class="produto-nome">${product.name}</span>
                <span class="produto-preco">R$ ${product.price.toFixed(2).replace('.', ',')}</span>
                <span class="produto-parcela">Até ${product.installments}x de R$ ${installmentValue.replace('.', ',')}</span>
            </div>
        </div>
    `;
}

function selectSize(productId, size) {
    document.querySelectorAll(`[data-product-id="${productId}"] .size-option`).forEach(el => {
        el.classList.remove('selected');
    });
    event.target.classList.add('selected');
    selectedSizes[productId] = size;
}

function addSelectedToCart(productId) {
    const selectedSize = selectedSizes[productId];

    if (!selectedSize) {
        showNotification('Por favor, selecione um tamanho primeiro!');
        return;
    }

    addToCart(productId, selectedSize);
    delete selectedSizes[productId];
    document.querySelectorAll(`[data-product-id="${productId}"] .size-option`).forEach(el => {
        el.classList.remove('selected');
    });
}

// FUNÇÕES DE FILTROS
function applyFilters() {
    const category = document.getElementById('categoryFilter').value;
    const size = document.getElementById('sizeFilter').value;
    const searchTerm = document.getElementById('searchProducts').value.toLowerCase().trim();

    // Filtrar produtos
    filteredProducts = allProducts.filter(product => {
        const matchesCategory = category === 'all' || product.category === category;
        const matchesSize = size === 'all' || product.sizes.includes(size);
        const matchesSearch = searchTerm === '' || product.name.toLowerCase().includes(searchTerm);
        
        return matchesCategory && matchesSize && matchesSearch;
    });
    
    renderProducts();
}

function searchProducts() {
    applyFilters();
}

function renderProducts() {
    const grid = document.getElementById('productsGrid');
    const countText = document.getElementById('productsCountText');
    
    if (filteredProducts.length === 0) {
        grid.innerHTML = `
            <div style="grid-column: 1 / -1; text-align: center; padding: 60px 20px;">
                <svg width="80" height="80" viewBox="0 0 24 24" fill="none" stroke="#ccc" stroke-width="1.5" style="margin-bottom: 20px;">
                    <circle cx="11" cy="11" r="8"/>
                    <path d="m21 21-4.35-4.35"/>
                </svg>
                <h3 style="font-family: 'Playfair Display', serif; font-size: 1.5rem; color: #370400; margin-bottom: 10px;">Nenhum produto encontrado</h3>
                <p style="color: #666; font-size: 14px;">Tente ajustar os filtros ou realizar uma nova busca.</p>
            </div>
        `;
        countText.textContent = '0';
    } else {
        grid.innerHTML = filteredProducts.map(product => createProductCard(product)).join('');
        countText.textContent = filteredProducts.length;
    }
}

// FUNÇÕES AUXILIARES
function showNotification(message) {
    const notification = document.createElement('div');
    notification.style.cssText = `
        position: fixed;
        top: 80px;
        right: 20px;
        background: #370400;
        color: white;
        padding: 15px 25px;
        border-radius: 4px;
        box-shadow: 0 4px 12px rgba(0,0,0,0.15);
        z-index: 10000;
        animation: slideIn 0.3s ease-out;
    `;
    notification.textContent = message;
    document.body.appendChild(notification);

    setTimeout(() => {
        notification.style.animation = 'slideOut 0.3s ease-out';
        setTimeout(() => notification.remove(), 300);
    }, 2000);
}

function toggleFilters() {
    document.getElementById('headerFilters').classList.toggle('hide');
}

function handleMenuToggleDisplay() {
    const btn = document.getElementById('menuToggle');
    const filters = document.getElementById('headerFilters');
    if (window.innerWidth <= 900) {
        btn.style.display = 'inline-flex';
    } else {
        btn.style.display = 'none';
        filters.classList.remove('hide');
    }
}

function initScrollBehavior() {
    const header = document.querySelector('.header');
    const footer = document.querySelector('.footer');
    let lastScroll = 0;
    let ticking = false;

    window.addEventListener('scroll', () => {
        if (!ticking) {
            window.requestAnimationFrame(() => {
                const currentScroll = window.scrollY;
                const footerRect = footer.getBoundingClientRect();
                const headerHeight = header.offsetHeight;

                if (footerRect.top <= headerHeight) {
                    header.classList.add('header-hide');
                } else {
                    if (currentScroll > lastScroll && currentScroll > 100) {
                        header.classList.add('header-hide');
                    } else if (currentScroll < lastScroll) {
                        header.classList.remove('header-hide');
                    }
                }

                lastScroll = currentScroll;
                ticking = false;
            });
            ticking = true;
        }
    });
}

function initFilterLinks() {
    document.querySelectorAll('.header-filters a').forEach(link => {
        link.addEventListener('click', function (e) {
            e.preventDefault();
            document.querySelectorAll('.header-filters a').forEach(l => l.classList.remove('active'));
            this.classList.add('active');
        });
    });
}

function initCategoriesModal() {
    const trigger = document.getElementById('categoriesMenuTrigger');
    const modal = document.getElementById('categoriesModal');
    if (!trigger || !modal) return;

    const closeBtn = document.getElementById('categoriesModalClose');

    if (closeBtn) {
        closeBtn.addEventListener('click', (e) => {
            e.stopPropagation();
            modal.classList.remove('active');
            document.body.classList.remove('modal-open');
        });
    }

    trigger.addEventListener('click', (e) => {
        e.preventDefault();
        e.stopPropagation();
        modal.classList.toggle('active');

        if (modal.classList.contains('active')) {
            document.body.classList.add('modal-open');
        } else {
            document.body.classList.remove('modal-open');
        }
    });

    document.addEventListener('click', (e) => {
        if (!trigger.contains(e.target) && !modal.contains(e.target)) {
            modal.classList.remove('active');
            document.body.classList.remove('modal-open');
        }
    });

    document.addEventListener('keydown', (e) => {
        if (e.key === 'Escape') {
            modal.classList.remove('active');
            document.body.classList.remove('modal-open');
        }
    });
}

// EVENT LISTENERS GLOBAIS
document.addEventListener('click', (e) => {
    const modal = document.getElementById('favoriteLoginModal');
    if (e.target === modal) {
        closeFavoriteModal();
    }
});

document.addEventListener('keydown', (e) => {
    if (e.key === 'Escape') {
        closeFavoriteModal();
    }
});

// INICIALIZAÇÃO
document.addEventListener('DOMContentLoaded', () => {
    renderProducts();
    updateFavoritesCount();
    updateCartCount();
    
    // Event listeners dos filtros
    document.getElementById('categoryFilter').addEventListener('change', applyFilters);
    document.getElementById('sizeFilter').addEventListener('change', applyFilters);
    document.getElementById('searchProducts').addEventListener('input', applyFilters);
    
    // Menu toggle
    document.getElementById('menuToggle').onclick = toggleFilters;
    window.addEventListener('resize', handleMenuToggleDisplay);
    handleMenuToggleDisplay();
    
    // Outros
    initScrollBehavior();
    initFilterLinks();
    initCategoriesModal();
});