// Função de busca do header que redireciona para /roupas
function initHeaderSearch() {
    const searchInput = document.getElementById('searchInput');
    const searchForm = document.querySelector('.search-form');
    const searchBtn = document.querySelector('.search-submit-btn');
    
    if (!searchInput || !searchForm) return;
    
    // Função para fazer a busca
    function realizarBusca(e) {
        if (e) e.preventDefault();
        
        const termoBusca = searchInput.value.trim();
        
        if (termoBusca) {
            // Redirecionar para /roupas com o termo de busca
            window.location.href = `/roupas?busca=${encodeURIComponent(termoBusca)}`;
        } else {
            // Se vazio, apenas redireciona para /roupas
            window.location.href = '/roupas';
        }
    }
    
    // Evento de submit do form
    searchForm.addEventListener('submit', realizarBusca);
    
    // Evento de clique no botão
    if (searchBtn) {
        searchBtn.addEventListener('click', realizarBusca);
    }
    
    // Evento de pressionar Enter no input
    searchInput.addEventListener('keypress', function(e) {
        if (e.key === 'Enter') {
            realizarBusca(e);
        }
    });
}

// Inicializar quando o DOM estiver pronto
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initHeaderSearch);
} else {
    initHeaderSearch();
}

