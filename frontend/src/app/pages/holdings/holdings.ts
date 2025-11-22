import { Component, OnInit, ChangeDetectorRef } from '@angular/core';
import { CommonModule } from '@angular/common';
import { TableModule } from 'primeng/table';
import { ButtonModule } from 'primeng/button';
import { DialogModule } from 'primeng/dialog';
import { InputTextModule } from 'primeng/inputtext';
import { InputNumberModule } from 'primeng/inputnumber';
// import { CalendarModule } from 'primeng/calendar'; //-- IGNORE ---
import { DatePickerModule } from 'primeng/datepicker'; // -- REPLACE CalendarModule WITH DatePickerModule ---
import { FormsModule } from '@angular/forms';
import { HoldingsService, Holding } from '../../services/holdings.service';
import { MessageService, ConfirmationService } from 'primeng/api';
import { ToastModule } from 'primeng/toast';
import { ConfirmDialogModule } from 'primeng/confirmdialog';

@Component({
  selector: 'app-holdings',
  standalone: true,
  imports: [
    CommonModule,
    TableModule,
    ButtonModule,
    DialogModule,
    InputTextModule,
    InputNumberModule,
    DatePickerModule,
    FormsModule,
    ToastModule,
    ConfirmDialogModule
  ],
  providers: [MessageService, ConfirmationService],
  templateUrl: './holdings.html',
  styleUrl: './holdings.scss',
})
export class Holdings implements OnInit {
  holdings: Holding[] = [];
  loading: boolean = true;
  displayDialog: boolean = false;
  holding: Holding = this.createEmptyHolding();
  submitted: boolean = false;

  constructor(
    private holdingsService: HoldingsService,
    private messageService: MessageService,
    private confirmationService: ConfirmationService,
    private cdr: ChangeDetectorRef
  ) {}

  ngOnInit() {
    this.loadHoldings();
  }

  loadHoldings() {
    this.loading = true;
    this.holdingsService.getHoldings().subscribe({
      next: (data) => {
        this.holdings = data;
        this.loading = false;
        this.cdr.detectChanges(); // Manually trigger change detection
      },
      error: (err) => {
        console.error('Failed to load holdings', err);
        this.messageService.add({ severity: 'error', summary: 'Error', detail: 'Failed to load holdings' });
        this.loading = false;
        this.cdr.detectChanges();
      }
    });
  }

  openNew() {
    this.holding = this.createEmptyHolding();
    this.submitted = false;
    this.displayDialog = true;
  }

  hideDialog() {
    this.displayDialog = false;
    this.submitted = false;
  }

  saveHolding() {
    this.submitted = true;

    if (this.holding.symbol.trim()) {
      if (this.holding.ID) {
        // Update
        this.holdingsService.updateHolding(this.holding.ID, this.holding).subscribe({
          next: (updatedHolding) => {
            const index = this.holdings.findIndex(h => h.ID === updatedHolding.ID);
            this.holdings[index] = updatedHolding;
            this.messageService.add({ severity: 'success', summary: 'Successful', detail: 'Holding Updated' });
            this.hideDialog();
          },
          error: (err) => {
            this.messageService.add({ severity: 'error', summary: 'Error', detail: 'Failed to update holding' });
          }
        });
      } else {
        // Create
        this.holdingsService.createHolding(this.holding).subscribe({
          next: (newHolding) => {
            this.holdings.push(newHolding);
            this.messageService.add({ severity: 'success', summary: 'Successful', detail: 'Holding Created' });
            this.hideDialog();
          },
          error: (err) => {
            this.messageService.add({ severity: 'error', summary: 'Error', detail: 'Failed to create holding' });
          }
        });
      }
    }
  }

  editHolding(holding: Holding) {
    this.holding = { ...holding };
    // Convert string date to Date object if necessary
    if (typeof this.holding.purchase_date === 'string') {
        this.holding.purchase_date = new Date(this.holding.purchase_date);
    }
    this.displayDialog = true;
  }

  deleteHolding(holding: Holding) {
    this.confirmationService.confirm({
      message: 'Are you sure you want to delete ' + holding.symbol + '?',
      header: 'Confirm',
      icon: 'pi pi-exclamation-triangle',
      accept: () => {
        if (holding.ID) {
            this.holdingsService.deleteHolding(holding.ID).subscribe({
                next: () => {
                    this.holdings = this.holdings.filter(val => val.ID !== holding.ID);
                    this.messageService.add({ severity: 'success', summary: 'Successful', detail: 'Holding Deleted' });
                },
                error: (err) => {
                    this.messageService.add({ severity: 'error', summary: 'Error', detail: 'Failed to delete holding' });
                }
            });
        }
      }
    });
  }

  createEmptyHolding(): Holding {
    return {
      symbol: '',
      quantity: 0,
      cost_basis: 0,
      purchase_date: new Date()
    };
  }
}
