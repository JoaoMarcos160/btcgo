package core

import (
	"btcgo/cmd/utils"
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

// Logger personalizado para exibir informações sobre o uso de variáveis de ambiente
var envLogger = log.New(os.Stdout, "[ENV] ", log.LstdFlags)

func RequestData() {
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	cpus := runtime.NumCPU()
	fmt.Printf("\nCPUs detectados: %s\n", green(cpus))

	// Configuração das propriedades da aplicação
	App.MaxWorkers = readCPUsForUse()
	if App.MaxWorkers == 0 {
		App.MaxWorkers = cpus
	}
	App.RangeNumber = promptRangeNumber()
	App.Carteira = fmt.Sprintf("%d", App.RangeNumber)
	App.Wallets.SetFindWallet(App.RangeNumber)
	App.Modo = promptMods(3)

	fmt.Printf("\n%s Configuração inicial:\n", yellow("✓"))
	fmt.Printf("  • CPUs em uso: %d/%d\n", App.MaxWorkers, cpus)
	fmt.Printf("  • Carteira selecionada: %s\n", App.Carteira)
	// Obter descrição do modo
	modoDescricao := ""
	switch App.Modo {
	case 1:
		modoDescricao = "Modo do início"
	case 2:
		modoDescricao = "Modo sequencial (chave do arquivo)"
	case 3:
		modoDescricao = "Modo Random"
	}
	fmt.Printf("  • Modo de operação: %d (%s)\n", App.Modo, modoDescricao)

	if App.Modo == 2 {
		App.DesdeInicio = false
		msSequencialouInicio := promptForEnvOrInput("START_MODE", "\n\nOpção 1: Deseja começar do inicio da busca (não efetivo) ou \nOpção 2: Escolher entre o range da carteira informada? \n\nPor favor numero entre 1 ou 2: ", "Número inválido. Escolha entre 1 ou 2.", 1, 2)

		if msSequencialouInicio == 1 {
			App.DesdeInicio = true
			fmt.Printf("  • Iniciando desde o começo: %t\n", App.DesdeInicio)
		} else {
			_, err := App.LastKey.GetLastKey(App.Carteira)
			if err != nil {
				// Solicitando a porcentagem do range da carteira como entrada
				App.StartPosPercent = promptForEnvOrFloat("START_PERCENT", "Informe a porcentagem do range da carteira entre 1 a 100: ")
				fmt.Printf("  • Porcentagem inicial do range: %.2f%%\n", App.StartPosPercent)
			}
		}
	} else if App.Modo == 3 {
		App.USEDB = promptUseDB(2)
		App.Keys.SetRecs(promptNumRecsRandom())
		fmt.Printf("  • Usando DB: %t\n", App.USEDB == 1)
		fmt.Printf("  • Registros por random: %d\n", App.Keys.NumRecsRandom)
	}

	fmt.Printf("%s Configuração concluída. Iniciando processamento...\n", yellow("✓"))
}

// Quantos CPUs gostaria de usar?
func readCPUsForUse() int {
	return promptForEnvOrInput("CPU_COUNT", "\n\nQuantos CPUs gostaria de usar?: ", "Número inválido.", 0, runtime.NumCPU())
}

// promptRangeNumber solicita ao usuário para selecionar um número de carteira
func promptRangeNumber() int {
	totalRanges := App.Ranges.Count()
	return promptForEnvOrInput("RANGE_NUMBER", fmt.Sprintf("\n\nEscolha a carteira (1 a %d): ", totalRanges), "Número inválido.", 1, totalRanges)
}

// promptMods solicita ao usuário a escolha de um modo
func promptMods(totalModos int) int {
	return promptForEnvOrInput("MODS", fmt.Sprintf("\n\nEscolha os modos que deseja de (1 a %d):\n1 - Modo do inicio\n2 - Modo sequencial (chave do arquivo)\n3 - Modo Random\n\nEscolha o modo: ", totalModos), "Modo inválido.", 1, totalModos)
}

// promptUseDB solicita ao usuário se deseja utilizar o banco de dados para controlar repetições
func promptUseDB(totalModos int) int {
	return promptForEnvOrInput("USE_DB", "\nUtiliza BaseDados para controlar repetições?\n1 - Modo Random com DB\n2 - Modo Random sem DB\n\nEscolha o modo: ", "Modo inválido.", 1, totalModos)
}

// promptNumRecsRandom solicita ao usuário o número de registros para o modo aleatório
func promptNumRecsRandom() int {
	return promptForEnvOrInput("NUM_RECS", "\nNúmero de registros por cada random (ex. 10000): ", "Número inválido.", 1, 0)
}

// promptForEnvOrInput tenta ler o valor de uma variável de ambiente ou solicita entrada do usuário
func promptForEnvOrInput(envVar, promptStr, errorStr string, min, max int) int {
	// Exibe a pergunta que seria feita ao usuário, para fins de log
	cleanPrompt := strings.ReplaceAll(promptStr, "\n", " ")
	fmt.Printf("\n[PROMPT] %s\n", cleanPrompt)

	if value, exists := os.LookupEnv(envVar); exists {
		if intValue, err := strconv.Atoi(value); err == nil && (max == 0 || (intValue >= min && intValue <= max)) {
			envLogger.Printf("Usando variável de ambiente %s=%s (em vez de entrada do usuário)", envVar, value)
			fmt.Printf("[RESPOSTA] %d (via variável de ambiente %s)\n", intValue, envVar)
			return intValue
		} else {
			// Erro apenas quando a variável existe mas é inválida
			log.Fatalf("[ERRO FATAL] Variável de ambiente %s='%s' é inválida: %v. O valor deve ser um número inteiro", envVar, value, err)
			// O código não chega aqui devido ao log.Fatalf
			return 0
		}
	} else {
		// Apenas log informativo quando a variável não existe - não é erro
		envLogger.Printf("Variável de ambiente %s não definida. Solicitando entrada do usuário.", envVar)
	}

	result := promptForIntInRange(promptStr, errorStr, min, max)
	fmt.Printf("[RESPOSTA] %d (via entrada do usuário)\n", result)
	return result
}

// promptForEnvOrFloat tenta ler um valor float de uma variável de ambiente ou solicita entrada do usuário
func promptForEnvOrFloat(envVar, promptStr string) float64 {
	// Exibe a pergunta que seria feita ao usuário, para fins de log
	cleanPrompt := strings.ReplaceAll(promptStr, "\n", " ")
	fmt.Printf("\n[PROMPT] %s\n", cleanPrompt)

	if value, exists := os.LookupEnv(envVar); exists {
		value = strings.Replace(value, ",", ".", -1)
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			envLogger.Printf("Usando variável de ambiente %s=%s (em vez de entrada do usuário)", envVar, value)
			fmt.Printf("[RESPOSTA] %.2f (via variável de ambiente %s)\n", floatValue, envVar)
			return floatValue
		} else {
			// Erro fatal apenas quando a variável existe mas é inválida
			log.Fatalf("[ERRO FATAL] Variável de ambiente %s='%s' é inválida: %v. O valor deve ser um número decimal", envVar, value, err)
			// O código não chega aqui devido ao log.Fatalf
			return 0.0
		}
	} else {
		// Apenas log informativo quando a variável não existe - não é erro
		envLogger.Printf("Variável de ambiente %s não definida. Solicitando entrada do usuário.", envVar)
	}

	var input string
	fmt.Print(promptStr)
	fmt.Scanln(&input)
	input = strings.Replace(input, ",", ".", -1)
	floatValue, err := strconv.ParseFloat(input, 64)
	if err != nil {
		envLogger.Printf("Entrada inválida: '%s'. Usando valor padrão 0.0", input)
		fmt.Printf("[RESPOSTA] 0.00 (valor padrão após erro na entrada)\n")
		return 0.0
	}
	fmt.Printf("[RESPOSTA] %.2f (via entrada do usuário)\n", floatValue)
	return floatValue
}

// promptForIntInRange solicita ao usuário a seleção de um número dentro de um intervalo específico.
// Será retornado um número que necessariamente atenda a (min <= X <= max) onde X foi o número escolhido pelo usuário.
func promptForIntInRange(requestStr, errorStr string, min, max int) int {
	charReadline := utils.GetEndLineChar()
	reader := bufio.NewReader(os.Stdin)
	attemptCount := 0
	for {
		attemptCount++
		if attemptCount > 1 {
			fmt.Printf("[TENTATIVA] %d - Solicitando novamente após entrada inválida\n", attemptCount)
		}

		fmt.Print(requestStr)
		input, _ := reader.ReadString(byte(charReadline))
		input = strings.TrimSpace(input)
		resposta, err := strconv.Atoi(input)

		if max == 0 {
			if err == nil && resposta >= min {
				return resposta
			}
		} else if err == nil && resposta >= min && resposta <= max {
			return resposta
		}

		if err != nil {
			fmt.Printf("[ERRO] Entrada '%s' não é um número válido\n", input)
		} else if max > 0 && (resposta < min || resposta > max) {
			fmt.Printf("[ERRO] Número %d está fora do intervalo permitido (%d a %d)\n", resposta, min, max)
		} else if resposta < min {
			fmt.Printf("[ERRO] Número %d é menor que o mínimo permitido (%d)\n", resposta, min)
		}

		fmt.Println(errorStr)
	}
}
