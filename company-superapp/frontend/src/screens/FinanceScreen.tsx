import { SafeAreaView, ScrollView, StyleSheet, Text, View } from 'react-native';
import SalaryWidget from '../components/SalaryWidget';
import { colors } from '../theme/colors';

export default function FinanceScreen() {
    return (
        <SafeAreaView style={styles.container}>
            <View style={styles.header}>
                <Text style={styles.title}>Финансы</Text>
            </View>
            <ScrollView contentContainerStyle={styles.content}>
                <SalaryWidget />
            </ScrollView>
        </SafeAreaView>
    );
}

const styles = StyleSheet.create({
    container: {
        flex: 1,
        backgroundColor: colors.background,
    },
    header: {
        paddingHorizontal: 16,
        paddingVertical: 20,
    },
    title: {
        color: colors.text,
        fontSize: 28,
        fontWeight: '700',
    },
    content: {
        paddingBottom: 20,
    },
});
