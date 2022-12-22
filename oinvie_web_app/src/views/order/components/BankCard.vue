<template>
    <div class="bank-card fixed-pop">
        <!--手勢下滑關閉-->
        <PullDownClose @close="cancel">
            <div class="fixed-pop-title">變更付款方式</div>
            <div class="mt-0 mr-4 ml-4 mb-4">
                <el-radio-group v-model="paymentRadio" @change="confirmClick()" class="radio-border-b bank-list">
                    <!--todo 價錢大於零樣式加red-->
                    <!--<el-radio label="1" class="bank-list" v-if="false">
                        <img src="@/assets/images/card-visa.svg" alt="visa">
                        <div class="bank">
                            <div class="name">台新國際商業銀行</div>
                            <div class="account">1234-5617-****-**36</div>
                        </div>
                        <span class="more">切換信用卡 <i class="el-icon-arrow-right"/></span>
                    </el-radio>-->
                    <el-radio label="1"><img src="@/assets/images/other_pay_cash.svg">現金</el-radio>
                    <el-radio label="2" v-if="paymentMethod['creditCard'] && !$public.isCordova()">
                        <img src="@/assets/images/other_pay_credit_card2.svg">信用卡
                    </el-radio>
                    <el-radio label="8" v-if="paymentMethod['creditCard'] && $public.isCordova()" :disabled="!paymentMethod['userPayToken']">
                        <img src="@/assets/images/other_pay_oinpay.svg">信用卡
                        <span v-show="paymentMethod['userPayToken']">({{creditCardNumber}})</span>
                        <el-button v-show="!paymentMethod['userPayToken']" @click="$router.push({name:'EmbedNewebpay'})" size="mini" round class="ml-5">
                            綁定信用卡付款
                        </el-button>
                    </el-radio>
                    <!--<el-radio label="3"><img src="@/assets/images/other_pay_linepay2.svg">LINE PAY</el-radio>
                    <el-radio label="4"><img src="@/assets/images/other_pay_easywallet2.svg">悠遊付</el-radio>
                    <el-radio label="5"><img src="@/assets/images/other_pay_pi2.svg">Pi錢包</el-radio>
                    <el-radio label="6"><img src="@/assets/images/other_pay_credit_card2.svg">支付寶</el-radio>
                    <el-radio label="7"><img src="@/assets/images/other_pay_wechatpay2.svg">微信支付</el-radio>-->
                </el-radio-group>
                <!--<el-button type="warning" @click="confirmClick()" class="bottom-btn">
                        確定
                </el-button>-->
            </div>
        </PullDownClose>
        <div class="fixed-pop-closed" @click="cancel()"/>
    </div>
</template>

<script>
    import PullDownClose from "../../../components/PullDownClose";
    export default {
        name: "BankCard",
        props: {
            paymentType:{},
            paymentMethod:{},
        },
        components: {
            PullDownClose,
        },
        data() {
            return {
                "paymentRadio": "",

                creditCardNumber: '',
            }
        },
        created() {
          this.paymentRadio = this.paymentType
        },
        mounted() {
          this.setupCreditCardData()
        },
        methods: {
            confirmClick() {
                if (this.paymentRadio === "") {
                    this.$public.showNotify("請選擇付款方式", false);
                    return false
                }
                this.$emit('changePaymentType', this.paymentRadio);
                //this.$parent.showBank = false
            },
            cancel() {
                //點空白處關閉
                this.$emit('changePaymentType', this.paymentRadio);
            },
            setupCreditCardData() {
                this.$http.fetchWithAuth`UserPayTokenByType${{
                    'tokenType': 1,
                }}
                ${json => {
                    if (json.status) {
                        this.creditCardNumber = json.data.creditCardNumber;
                    } else {
                        this.creditCardNumber = "000000******0000";
                    }
                }}`;
            },
        }
    }
</script>

