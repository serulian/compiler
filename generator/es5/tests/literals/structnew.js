$module('structnew', function () {
  var $static = this;
  this.$class('SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (SomeField) {
      var instance = new $static();
      instance.SomeField = SomeField;
      instance.anotherField = $t.box(false, $g.____testlib.basictypes.Boolean);
      return $promise.resolve(instance);
    };
    $instance.set$AnotherField = function (val) {
      var $this = this;
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        $this.anotherField = val;
        $resolve();
        return;
      };
      return $promise.new($continue);
    };
    $instance.AnotherField = $t.property(function () {
      var $this = this;
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        $resolve($this.anotherField);
        return;
      };
      return $promise.new($continue);
    });
    this.$typesig = function () {
      var computed = $t.createtypesig(['AnotherField', 3, $g.____testlib.basictypes.Boolean.$typeref()], ['new', 1, $g.____testlib.basictypes.Function($g.structnew.SomeClass).$typeref()]);
      this.$typesig = function () {
        return computed;
      };
      return computed;
    };
  });

  $static.TEST = function () {
    var $result;
    var $temp0;
    var sc;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.structnew.SomeClass.new($t.box(2, $g.____testlib.basictypes.Integer)).then(function ($result0) {
              $temp0 = $result0;
              return $temp0.set$AnotherField($t.box(true, $g.____testlib.basictypes.Boolean)).then(function ($result1) {
                $result = ($temp0, $result1, $temp0);
                $current = 1;
                $continue($resolve, $reject);
                return;
              });
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 1:
            sc = $result;
            $g.____testlib.basictypes.Integer.$equals(sc.SomeField, $t.box(2, $g.____testlib.basictypes.Integer)).then(function ($result1) {
              return $promise.resolve($t.unbox($result1)).then(function ($result0) {
                $result = $t.box($result0 && $t.unbox(sc.anotherField), $g.____testlib.basictypes.Boolean);
                $current = 2;
                $continue($resolve, $reject);
                return;
              });
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 2:
            $resolve($result);
            return;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  };
});
