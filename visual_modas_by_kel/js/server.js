const express = require('express');
const { MercadoPagoConfig, Preference, Payment } = require('mercadopago');
const cors = require('cors');

const app = express();
app.use(express.json());
app.use(cors());

// Configure o Mercado Pago com suas credenciais (versão 2.x)
const client = new MercadoPagoConfig({
    accessToken: 'TEST-2974489993984774-100215-421b7b763b7b0cfe09fcdc4125783e43-1162798817',
    options: { timeout: 5000 }
});

// ENDPOINT: Criar pagamento com Mercado Pago
app.post('/api/criar-pagamento-mp', async (req, res) => {
    try {
        const { cliente, endereco, items, total } = req.body;

        // Formata os itens para o Mercado Pago
        const itemsMP = items.map(item => ({
            title: item.name,
            description: `Tamanho: ${item.size}`,
            quantity: item.quantity,
            unit_price: item.price,
            currency_id: 'BRL'
        }));

        // Cria a preferência de pagamento
        const preference = new Preference(client);
        
        const preferenceData = {
            items: itemsMP,
            payer: {
                name: cliente.nome,
                email: cliente.email,
                identification: {
                    type: 'CPF',
                    number: cliente.cpf
                },
                phone: {
                    area_code: cliente.telefone.substring(0, 2),
                    number: cliente.telefone.substring(2)
                },
                address: {
                    zip_code: endereco.cep,
                    street_name: endereco.endereco,
                    street_number: parseInt(endereco.numero)
                }
            },
            back_urls: {
                success: 'http://localhost/checkout-sucesso.html',
                failure: 'http://localhost/checkout-falha.html',
                pending: 'http://localhost/checkout-pendente.html'
            },
            auto_return: 'approved',
            statement_descriptor: 'VISUAL MODAS BY KEL',
            external_reference: `PEDIDO-${Date.now()}`
        };

        const response = await preference.create({ body: preferenceData });

        res.json({
            success: true,
            preferenceId: response.id,
            initPoint: response.init_point,
            transactionId: response.external_reference
        });

    } catch (error) {
        console.error('Erro ao criar pagamento MP:', error);
        res.status(500).json({
            success: false,
            message: 'Erro ao processar pagamento',
            error: error.message
        });
    }
});

// ENDPOINT: Criar pagamento PIX
app.post('/api/criar-pagamento-pix', async (req, res) => {
    try {
        const { cliente, endereco, items, total } = req.body;

        const payment = new Payment(client);
        
        const paymentData = {
            transaction_amount: total,
            description: 'Pedido Visual Modas By Kel',
            payment_method_id: 'pix',
            payer: {
                email: cliente.email,
                first_name: cliente.nome.split(' ')[0],
                last_name: cliente.nome.split(' ').slice(1).join(' ') || cliente.nome.split(' ')[0],
                identification: {
                    type: 'CPF',
                    number: cliente.cpf
                }
            }
        };

        const response = await payment.create({ body: paymentData });

        res.json({
            success: true,
            qrCode: response.point_of_interaction.transaction_data.qr_code,
            qrCodeBase64: response.point_of_interaction.transaction_data.qr_code_base64,
            transactionId: response.id
        });

    } catch (error) {
        console.error('Erro ao criar pagamento PIX:', error);
        res.status(500).json({
            success: false,
            message: 'Erro ao processar pagamento PIX',
            error: error.message
        });
    }
});

// ENDPOINT: Webhook do Mercado Pago
app.post('/api/webhook-mp', async (req, res) => {
    try {
        const { type, data } = req.body;

        if (type === 'payment') {
            const payment = new Payment(client);
            const paymentInfo = await payment.get({ id: data.id });
            
            console.log('Status do pagamento:', paymentInfo.status);
            console.log('ID do pedido:', paymentInfo.external_reference);
            
            // Aqui você processaria o status do pagamento
            // if (paymentInfo.status === 'approved') {
            //     // Pagamento aprovado - atualizar banco de dados
            // }
        }

        res.sendStatus(200);
    } catch (error) {
        console.error('Erro no webhook:', error);
        res.sendStatus(500);
    }
});

// Rota de teste
app.get('/api/test', (req, res) => {
    res.json({ message: 'Servidor funcionando!' });
});

const PORT = process.env.PORT || 3000;
app.listen(PORT, () => {
    console.log(`Servidor rodando na porta ${PORT}`);
    console.log(`Acesse: http://localhost:${PORT}/api/test`);
});