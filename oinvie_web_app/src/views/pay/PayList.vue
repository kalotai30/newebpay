<template>
    <div>
        <navigation-bar Back/>
        <div class="main container">
            <div class="d-flex pay-list">
                <img src="../../assets/images/pay_card1.svg" class="fixed-w-70 ml-3">
                <div class="flex-grow-1">
                    <h3>開啟信用卡條碼</h3>
                    <p>使用付款條碼進行線上付款</p>
                </div>
            </div>
            <div class="shadow-line mt-0 mb-0"></div>
            <!-- <div class="d-flex pay-list" @click="$router.push({name:'CreditCardList'})"> -->
            <div class="d-flex pay-list" @click="$router.push({name:'EmbedNewebpay'})">
                <img src="../../assets/images/pay_card2.svg" class="fixed-w-50">
                <div class="flex-grow-1">
                    <h3>我的信用卡</h3>
                    <p>新增/編輯 您所使用的信用卡</p>
                </div>
            </div>
            <!--
            <el-divider class="m-0"/>
            <div class="d-flex pay-list" @click="$router.push({name:'PayEdit'})">
            <div class="d-flex pay-list">
                <img src="../../assets/images/pay_card2.svg" class="fixed-w-50">
                <div class="flex-grow-1">
                    <h3>設定</h3>
                    <p>Touch ID、更換支付密碼</p>
                </div>
            </div>
            -->
        </div>
    </div>
</template>

<script>

    export default {
        name: "PayList",
        data() {
            return {
                
            }
        },
        created() {
            
        },
        mounted() {
            
        },
        methods: {
             
        }

    }
</script>

<style lang="less" scoped src="../../less/pay.less"/>
