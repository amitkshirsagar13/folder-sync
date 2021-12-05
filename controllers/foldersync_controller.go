/*
Copyright 2021 Amit Kshirsagar.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	foldersv1 "amn.siemens.com/m/v2/api/v1"
	v1 "amn.siemens.com/m/v2/api/v1"
)

// FolderSyncReconciler reconciles a FolderSync object
type FolderSyncReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=folders.operators.amn.siemens.com,resources=foldersyncs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=folders.operators.amn.siemens.com,resources=foldersyncs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=folders.operators.amn.siemens.com,resources=foldersyncs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the FolderSync object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *FolderSyncReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// TODO(user): your logic here

	log.Info("Running in reconcile!!!")

	fs := &v1.FolderSync{}
	if err := r.Client.Get(ctx, req.NamespacedName, fs); err != nil {
		// add some debug information if it's not a NotFound error
		// if !k8serr.IsNotFound(err) {
		// 	log.Error(err, "unable to fetch VmGroup")
		// }
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		log.Info("error in reconcile!!!")
		r.Client.Status().Update(context.TODO(), fs)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	current := int32(0)

	currentWd, error := os.Getwd()

	folderExists, error := exists(fs.Spec.FolderName)
	if !folderExists {
		os.Mkdir(fs.Spec.FolderName, 0777)
	}

	folderExists, error = exists(fs.Spec.FolderName)
	fileLen := 0
	if folderExists && current != 5 {
		files, _ := ioutil.ReadDir(fmt.Sprint(currentWd, "/", fs.Spec.FolderName))
		current = int32(len(files))
		for _, f := range files {
			fmt.Println(f.Name())
		}
		subFolderName := fmt.Sprint("")
		if current < 5 {
			current = current + 1
			subFolderName = fmt.Sprint(currentWd, "/", fs.Spec.FolderName, "/", current)
			os.MkdirAll(subFolderName, 0777)
		} else {
			if current != 5 {
				for _, f := range files {
					subFolderName = fmt.Sprint(currentWd, "/", fs.Spec.FolderName, "/", f.Name())
				}
				error = os.Remove(subFolderName)
			}
		}
		fileLen = len(files)
	}
	message := fmt.Sprint("Status Current ", currentWd, "|", folderExists, "|", fileLen, "|", current)
	log.Info(message)

	status := createStatus(fs.Spec.FolderName, folderExists, &current, fs.Spec.SubFolderCount)
	fs.Status = status
	log.Info("status in reconcile!!!")

	error = r.Client.Status().Update(context.TODO(), fs)
	return ctrl.Result{}, error
}

// SetupWithManager sets up the controller with the Manager.
func (r *FolderSyncReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&foldersv1.FolderSync{}).
		Complete(r)
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func createStatus(folderName string, folderExists bool, current *int32, desired int32) v1.FolderSyncStatus {

	status := v1.FolderSyncStatus{
		FolderName:            folderName,
		FolderNameExists:      folderExists,
		CurrentSubFolderCount: current,
		DesiredSubFolderCount: desired,
	}
	return status
}
