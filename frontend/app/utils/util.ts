export function formatCurrency(value: number | null | undefined): string {
    if (value === null || value === undefined) {
        return '0,00 €';
    }

    return new Intl.NumberFormat('de-DE', {
        style: 'currency',
        currency: 'EUR',
    }).format(value);
}